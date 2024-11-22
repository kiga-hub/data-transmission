package upgrade

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kiga-hub/arc/logging"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Req -
type Req struct {
	ProjectName string `json:"project_name"`
	RemoteIP    string `json:"remote_ip"` // "ip:port"
	User        string `json:"user"`
	Password    string `json:"password"`
}

// Client -
type Client struct {
	client *ssh.Client
	config *Config
	logger logging.ILogger
	// out      io.WriteCloser
	// redirect *log.Logger
}

// RemoteTarget -
type RemoteTarget struct {
	client   *ssh.Client
	logger   logging.ILogger
	out      io.WriteCloser
	redirect *log.Logger
}

// New -
func New(opts ...Option) (*Client, error) {
	srv := loadOptions(opts...)

	// conf := GetConfig()
	// fmt.Println("conf: ", conf)
	// config := &ssh.ClientConfig{
	// 	User: conf.User,
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.Password(conf.Password),
	// 	},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }

	// client, err := ssh.Dial("tcp", conf.SourceIP, config)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to dial SSH: %v", err)
	// }
	srv.client = nil
	return srv, nil
}

// NewRemoteTarget -
func NewRemoteTarget(param *Req, logger logging.ILogger) (*RemoteTarget, error) {
	config := &ssh.ClientConfig{
		User: param.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(param.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", param.RemoteIP, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %v", err)
	}

	date := time.Now().Format("20060102150405")
	logDir := filepath.Join("./upgrade_log", param.ProjectName, date)

	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, "upgrade.log")
	out, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	redirect := log.New(out, "", log.LstdFlags)

	remoteClient := &RemoteTarget{
		logger:   logger,
		client:   client,
		out:      out,
		redirect: redirect,
	}
	return remoteClient, nil
}

// Close -
func (s *RemoteTarget) Close() {
	s.out.Close()
	s.client.Close()
}

// Write -
func (s *RemoteTarget) Write(p []byte) (n int, err error) {
	fmt.Printf("%s", string(p))
	return s.out.Write(p)
}

// RunRemoteCommand -
func (s *RemoteTarget) RunRemoteCommand(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v, output: %s", err, output)
	}

	return string(output), nil
}

// DetermineRemoteShmSize shm-size based on memory using SSHClient
func (s *RemoteTarget) DetermineRemoteShmSize() (string, error) {
	memGB, err := s.GetRemoteTotalMemory()
	if err != nil {
		s.redirect.Printf("Error getting total memory: %v\n", err)
		return "", err
	}
	if memGB < 250 {
		s.redirect.Printf("Memory is less than 250GB. No specific shm-size recommendation.\n")
		return "--shm-size=32G", nil
	} else if memGB > 500 {
		s.redirect.Printf("Memory is greater than 500GB. No specific shm-size recommendation.\n")
		return "--shm-size=64G", nil
	}
	s.redirect.Printf("Memory is between 250GB and 500GB. No specific shm-size recommendation.\n")
	return "", nil
}

// CheckRemoteNvidiaGPU - Check if remote server has Nvidia GPU
func (s *RemoteTarget) CheckRemoteNvidiaGPU() (string, string, bool) {
	cmd := "nvidia-smi"
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		return "", output, false
	}
	if strings.Contains(output, "GPU") {
		return "--gpus all", output, true
	}
	return "", output, false
}

// GetRemoteTotalMemory - Get total memory of remote server
func (s *RemoteTarget) GetRemoteTotalMemory() (int, error) {
	cmd := "grep MemTotal /proc/meminfo"
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		return 0, err
	}
	memInfo := strings.Fields(output)
	memKB, err := strconv.Atoi(memInfo[1])
	if err != nil {
		return 0, err
	}
	memGB := memKB / 1024 / 1024
	return memGB, nil
}

// CheckRemoteRunningContainer Check if a container with the specified image is already running using SSHClient
func (s *RemoteTarget) CheckRemoteRunningContainer(image string) (string, error) {
	cmd := fmt.Sprintf("docker ps -q --filter ancestor=%s", image)
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Error checking running container: %v\n", err)
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// StopAndRemoveRemoteContainer -
func (s *RemoteTarget) StopAndRemoveRemoteContainer(containerID string) error {
	// Stop the container
	stopCmd := fmt.Sprintf("docker stop %s", containerID)
	stopOutput, err := s.RunRemoteCommand(stopCmd)
	if err != nil {
		s.redirect.Printf("Error stopping container: %v\n", err)
		return fmt.Errorf("failed to stop container: %v, output: %s", err, stopOutput)
	}
	s.redirect.Printf("Container stopped: %s\n", stopOutput)

	// Remove the container
	rmCmd := fmt.Sprintf("docker rm %s", containerID)
	rmOutput, err := s.RunRemoteCommand(rmCmd)
	if err != nil {
		s.redirect.Printf("Error removing container: %v\n", err)
		return fmt.Errorf("failed to remove container: %v, output: %s", err, rmOutput)
	}
	s.redirect.Printf("Container removed: %s\n", rmOutput)

	return nil
}

// MakeRemoteContainerCommand -
func (s *RemoteTarget) MakeRemoteContainerCommand(sourceConfig SourceConfig) string {
	gpuOption := ""
	gpuDetected := false
	cmdOutput := ""

	gpuOption, cmdOutput, gpuDetected = s.CheckRemoteNvidiaGPU()
	if gpuDetected {
		s.redirect.Printf("NVIDIA GPU is installed: %s\n", gpuOption)
		s.logger.Info("NVIDIA GPU is installed:", gpuOption)
		s.logger.Info("nvidia-smi Output:\n", cmdOutput)
		s.redirect.Printf("nvidia-smi Output:\n%s\n", cmdOutput)
	} else {
		s.logger.Info("NVIDIA GPU is not installed. Remove --gpus all option.")
		s.redirect.Printf("NVIDIA GPU is not installed. Remove --gpus all option.\n")
	}

	shmSize, err := s.DetermineRemoteShmSize()
	if err != nil {
		s.redirect.Printf("Error determining shm-size: %v\n", err)
		s.logger.Info("Error determining shm-size:", err)
	} else if shmSize != "" {
		s.redirect.Printf("Recommended shm-size: %s\n", shmSize)
		s.logger.Info("Recommended shm-size:", shmSize)
	} else {
		s.redirect.Printf("Memory is between 250GB and 500GB. No specific shm-size recommendation.\n")
		s.logger.Info("Memory is between 250GB and 500GB. No specific shm-size recommendation.")
	}
	s.redirect.Printf("gpuSupport: %s\n", gpuOption)
	s.logger.Info("gpuSupport:", gpuOption)
	s.redirect.Printf("shmSize: %s\n", shmSize)
	s.logger.Info("shmSize:", shmSize)

	dockerCmd := "docker run --privileged"
	if gpuOption != "" {
		dockerCmd += " " + gpuOption
	}
	if shmSize != "" {
		dockerCmd += " " + shmSize
	}

	// TODO container image
	containerImage := ""
	dockerCmd += " -it --restart=unless-stopped -d -p1234:1234"
	dockerCmd += " " + containerImage

	s.redirect.Printf("Docker command: %s\n", dockerCmd)

	return dockerCmd
}

// StartRemoteContainer -
func (s *RemoteTarget) StartRemoteContainer(sourceConfig SourceConfig) error {
	s.redirect.Println("Start Remote Container: ", sourceConfig)
	dockerCmd := s.MakeRemoteContainerCommand(sourceConfig)
	// containerID, err := s.CheckRemoteRunningContainer(sourceConfig.ContainerImage)
	// if err != nil {
	// 	s.redirect.Printf("Error checking running container: %v\n", err)
	// 	fmt.Println("Error checking running container:", err)
	// 	return err
	// }

	// if containerID != "" {
	// 	s.redirect.Printf("container running with image: %s\n", sourceConfig.ContainerImage)
	// 	s.redirect.Printf("Stopping and removing existing container: %s\n", containerID)

	// 	fmt.Println("Stopping and removing existing container:", containerID)
	// 	err = s.StopAndRemoveRemoteContainer(containerID)
	// 	if err != nil {
	// 		s.redirect.Printf("Error stopping/removing container: %v\n", err)
	// 		fmt.Println("Error stopping/removing container:", err)
	// 		return err
	// 	}
	// }

	fmt.Println("\nDocker command\n", dockerCmd)
	s.redirect.Printf("Docker command: %s\n", dockerCmd)

	// Execute the Docker command
	output, err := s.RunRemoteCommand(fmt.Sprintf("sh -c '%s'", dockerCmd))
	if err != nil {
		fmt.Println("Error running Docker command:", err)
		fmt.Println(output)
		return err
	}
	s.redirect.Printf("Docker command output: %s\n", output)

	// Get the container ID
	containerID := strings.TrimSpace(output)
	s.redirect.Printf("Container ID: %s\n", containerID)
	fmt.Println("Container ID:", containerID)

	// Monitor the container logs until "Started Metrics Service" is found
	logCmd := fmt.Sprintf("docker logs %s | head -n 20", containerID)
	logOutput, err := s.RunRemoteCommand(logCmd)
	if err != nil {
		s.redirect.Printf("Error getting container logs: %v\n", err)
		fmt.Println("Error getting container logs:", err)
		return err
	}

	s.redirect.Printf("Container logs: %s\n", logOutput)
	fmt.Println(logOutput)
	// Sleep for a short duration before checking logs again
	time.Sleep(1 * time.Second)

	return nil
}

// CheckRemoteDirectory - checks if a directory exists on a remote server, and if not, creates it.
func (s *RemoteTarget) CheckRemoteDirectory(directory string) error {
	// Execute the command to check if the directory exists
	checkCommand := fmt.Sprintf("if [ -d %s ]; then echo 'exists'; else echo 'not exists'; fi", directory)
	output, err := s.RunRemoteCommand(checkCommand)
	if err != nil {
		s.redirect.Printf("Failed to check directory %s: %v\n", checkCommand, err)
		fmt.Printf("Failed to create backup directory %s: %v\n", checkCommand, err)
		return err
	}
	s.redirect.Printf("Directory check output: %s\n", output)

	// If the directory does not exist, create it
	if string(output) == "not exists\n" {
		createCommand := fmt.Sprintf("mkdir -p %s", directory)
		output, err := s.RunRemoteCommand(createCommand)
		if err != nil {
			s.redirect.Printf("Failed to create backup directory %s: %v\n", createCommand, err)
			fmt.Printf("Failed to create backup directory %s: %v\n", createCommand, err)
			return err
		}
		s.redirect.Printf("%s", output)
		s.redirect.Printf("Directory created\n")
		fmt.Println("Directory created")
	} else {
		s.redirect.Printf("Directory exists\n")
		fmt.Println("Directory exists")
	}

	return nil
}

// ImportSQLToRemoteDatabase -
func (s *RemoteTarget) ImportSQLToRemoteDatabase(containerID, scriptName string, mysqlConfig []ModelConfig) error {
	s.redirect.Printf("Start importing SQL script into database: %s\n", scriptName)
	fmt.Println("Start importing SQL script into database: ", scriptName)
	user := ""
	password := ""
	databaseName := ""
	for _, config := range mysqlConfig {
		if config.Key == "user" {
			user = config.Value
		}
		if config.Key == "password" {
			password = config.Value
		}
		if config.Key == "database" {
			databaseName = config.Value
		}
	}

	if user != "" && password != "" && databaseName != "" {
		fmt.Println("user: ", user)
		fmt.Println("password: ", password)
		fmt.Println("database: ", databaseName)

		dockerCmdStr := fmt.Sprintf("docker exec -i %s mysql -u %s -p%s %s < %s", containerID, user, password, databaseName, scriptName)
		output, err := s.RunRemoteCommand(dockerCmdStr)
		if err != nil {
			s.redirect.Printf("Error importing SQL script: %v\n", err)
			return fmt.Errorf("error running command: %v\noutput: %s", err, output)
		}
		s.redirect.Printf("%s\n", output)
		s.redirect.Printf("SQL script imported successfully\n")
		fmt.Println(output)
	}
	return nil
}

// RemoveRemoteFile -
func (s *RemoteTarget) RemoveRemoteFile(filePath string) error {
	// Check if the file exists
	checkCmd := fmt.Sprintf("ls %s", filePath)
	output, err := s.RunRemoteCommand(checkCmd)
	if err != nil {
		s.redirect.Printf("Failed to check if file exists: %v", err)
		return fmt.Errorf("file does not exist: %v", err)
	}
	s.redirect.Printf("%s", output)

	// Remove the file
	removeCmd := fmt.Sprintf("rm %s", filePath)
	output, err = s.RunRemoteCommand(removeCmd)
	if err != nil {
		s.redirect.Printf("Failed to delete remote file: %v", err)
		return fmt.Errorf("failed to delete remote file: %v", err)
	}
	s.redirect.Printf("%s", output)
	fmt.Println("File successfully deleted:", filePath)
	return nil
}

// BackupRemoteFile ackup a file or directory to a specified directory
func (s *RemoteTarget) BackupRemoteFile(srcPath, backupDir string) error {
	// get current time to avoid duplicate file names
	timeStamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s_%s", timeStamp, filepath.Base(srcPath)))

	// move file or directory to backup directory
	cmd := fmt.Sprintf("mv %s %s", srcPath, backupPath)
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Failed to backup %s to %s: %v", srcPath, backupPath, err)
		return fmt.Errorf("failed to backup %s to %s: %v", srcPath, backupPath, err)
	}
	s.redirect.Printf("Backed up %s to %s\n", srcPath, backupPath)
	s.redirect.Printf("Backup output: %s\n", output)
	fmt.Printf("Backed up %s to %s\n", srcPath, backupPath)
	return nil
}

// ExtractRemoteTarGz nzip the tar.gz file to the destination directory
func (s *RemoteTarget) ExtractRemoteTarGz(tarGzPath, destDir, backupDir string) error {
	// create the command to extract the tar.gz file to the destination directory
	// cmd := fmt.Sprintf("mkdir -p %s && tar -xzf %s -C %s", destDir, tarGzPath, destDir)
	// output, err := s.RunRemoteCommand(cmd)
	// if err != nil {
	// 	s.redirect.Printf("Failed to extract tar.gz file %s to %s: %v", tarGzPath, destDir, err)
	// 	return fmt.Errorf("failed to extract tar.gz file %s to %s: %v", tarGzPath, destDir, err)
	// }
	// s.redirect.Printf("Tar.gz contents: %s\n", output)

	// get the list of files in the tar.gz file
	cmd := fmt.Sprintf("tar -tzf %s", tarGzPath)
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Failed to list tar.gz contents %s: %v", tarGzPath, err)
		return fmt.Errorf("failed to list tar.gz contents %s: %v", tarGzPath, err)
	}
	s.redirect.Printf("Tar.gz contents: %s\n", output)

	// handle each file
	files := strings.Split(output, "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		destPath := filepath.Join(destDir, file)

		// check if the destination file or directory already exists
		if _, err := s.RunRemoteCommand(fmt.Sprintf("stat %s", destPath)); err == nil {
			// if the destination file or directory already exists, back it up
			err := s.BackupRemoteFile(destPath, backupDir)
			if err != nil {
				return err
			}

			// delete existing file or directory
			cmd := fmt.Sprintf("rm -rf %s", destPath)
			output, err = s.RunRemoteCommand(cmd)
			if err != nil {
				s.redirect.Printf("Failed to remove existing file or directory %s: %v", destPath, err)
				return fmt.Errorf("failed to remove existing file or directory %s: %v", destPath, err)
			}
			s.redirect.Printf("Removed existing file or directory: %s\n", destPath)
			s.redirect.Printf("%s", output)
			fmt.Printf("Removed existing file or directory: %s\n", destPath)
		}
	}

	// unzip the tar.gz file to the destination directory
	cmd = fmt.Sprintf("tar -xzf %s -C %s", tarGzPath, destDir)
	output, err = s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Failed to extract tar.gz file %s to %s: %v", tarGzPath, destDir, err)
		return fmt.Errorf("failed to extract tar.gz file %s to %s: %v", tarGzPath, destDir, err)
	}
	s.redirect.Printf("Extracted tar.gz file: %s to %s\n", tarGzPath, destDir)
	s.redirect.Printf("%s", output)

	fmt.Printf("Extracted tar.gz file: %s to %s\n", tarGzPath, destDir)
	return nil
}

// DownloadTarball2Remote - downlaod a tarball from a remote URL to the target server using wget
func (s *RemoteTarget) DownloadTarball2Remote(url string, targetDir string) error {
	// Create the target directory if it doesn't exist
	mkdirCmd := fmt.Sprintf("mkdir -p %s", targetDir)
	_, err := s.RunRemoteCommand(mkdirCmd)
	if err != nil {
		s.redirect.Printf("Error creating target directory: %v\n", err)
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// Download the tarball using wget
	wgetCmd := fmt.Sprintf("wget -P %s %s", targetDir, url)
	output, err := s.RunRemoteCommand(wgetCmd)
	if err != nil {
		s.redirect.Printf("Error downloading tarball: %v\n", err)
		return fmt.Errorf("failed to download tarball: %v, output: %s", err, output)
	}

	s.redirect.Printf("Successfully downloaded tarball %s to %s\n", url, targetDir)
	return nil
}

// TransferTarball2Remote - Transfer a tarball from resource IP to target IP
func (s *RemoteTarget) TransferTarball2Remote(resourceIP, targetIP, remoteFilePath, localFilePath string, TargetClient *Client) error {
	// Create a new SFTP client for resource IP
	sftpClientResource, err := sftp.NewClient(s.client)
	if err != nil {
		s.redirect.Printf("Failed to create SFTP client at resource IP: %v", err)
		log.Fatalf("Failed to create SFTP client at resource IP: %v", err)
	}
	defer sftpClientResource.Close()
	// Open the remote file at resource IP
	remoteFile, err := sftpClientResource.Open(remoteFilePath)
	if err != nil {
		s.redirect.Printf("Failed to open remote file at resource IP: %v", err)
		log.Fatalf("Failed to open remote file at resource IP: %v", err)
	}
	defer remoteFile.Close()

	// Create a new SFTP client for target IP
	sftpClientTarget, err := sftp.NewClient(TargetClient.client)
	if err != nil {
		s.redirect.Printf("Failed to create SFTP client at target IP: %v", err)
		log.Fatalf("Failed to create SFTP client at target IP: %v", err)
	}
	defer sftpClientTarget.Close()
	// Create the remote file at target IP
	targetFile, err := sftpClientTarget.Create(localFilePath)
	if err != nil {
		s.redirect.Printf("Failed to create remote file at target IP: %v", err)
		log.Fatalf("Failed to create remote file at target IP: %v", err)
	}
	defer targetFile.Close()

	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		s.redirect.Printf("Failed to stat remote file: %v", err)
		return fmt.Errorf("failed to stat remote file: %v", err)
	}

	// Copy the remote file from resource IP to target IP with progress logging
	buf := make([]byte, 1024)
	totalBytes := int64(0)
	startTime := time.Now()
	lastLoggedPercentage := float64(0)
	percentageThreshold := float64(10) // Log every 10%
	for {
		n, err := remoteFile.Read(buf)
		if err != nil && err != io.EOF {
			s.redirect.Printf("Failed to read from remote file at resource IP: %v", err)
			return fmt.Errorf("failed to read from remote file at resource IP: %v", err)
		}
		if n == 0 {
			break
		}
		if _, err := targetFile.Write(buf[:n]); err != nil {
			s.redirect.Printf("Failed to write to remote file at target IP: %v", err)
			return fmt.Errorf("failed to write to remote file at target IP: %v", err)
		}
		totalBytes += int64(n)
		elapsedTime := time.Since(startTime).Seconds()
		transferRate := float64(totalBytes) / elapsedTime
		percentage := float64(totalBytes) / float64(remoteFileInfo.Size()) * 100
		// fmt.Printf("\rTransferred: %d bytes (%.2f%%) at %.2f bytes/sec", totalBytes, percentage, transferRate)
		// s.redirect.Printf("Transferred: %d bytes (%.2f%%) at %.2f bytes/sec", totalBytes, percentage, transferRate)
		if percentage >= lastLoggedPercentage+percentageThreshold || totalBytes == remoteFileInfo.Size() {
			s.redirect.Printf("Transferred: %d bytes (%.2f%%) at %.2f bytes/sec", totalBytes, percentage, transferRate)
			lastLoggedPercentage = percentage
		}
	}
	s.redirect.Printf("File transferred successfully")
	fmt.Println("File transferred successfully")
	return nil
}

// CreateRemoteDirectories create target directory and backup directory
func (s *RemoteTarget) CreateRemoteDirectories(destDir, backupDir string) error {
	cmd := fmt.Sprintf("mkdir -p %s", destDir)
	output, err := s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Failed to create target directory %s: %v", destDir, err)
		fmt.Printf("Failed to create target directory %s: %v\n", destDir, err)
		return err
	}
	s.redirect.Printf("%s", output)

	// create the backup directory
	cmd = fmt.Sprintf("mkdir -p %s", backupDir)
	output, err = s.RunRemoteCommand(cmd)
	if err != nil {
		s.redirect.Printf("Failed to create backup directory %s: %v", backupDir, err)
		fmt.Printf("Failed to create backup directory %s: %v\n", backupDir, err)
		return err
	}
	s.redirect.Printf("%s", output)

	return nil
}

// TransmitRemoteSource -
func (s *RemoteTarget) TransmitRemoteSource(packagePath, destDir, backupDir string) error {
	err := s.CreateRemoteDirectories(destDir, backupDir)
	if err != nil {
		return err
	}

	err = s.ExtractRemoteTarGz(packagePath, destDir, backupDir)
	if err != nil {
		return err
	}

	fmt.Println("Extraction completed successfully.")
	s.redirect.Printf("Extraction completed successfully.")
	return nil
}
