package upgrade

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/kiga-hub/data-transmission/pkg/utils"
)

// UpdateSourceList -
func (s *Client) UpdateSourceList(url string) ([]*SourceConfig, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("GetSourceList cost time: ", time.Since(start))
	}()

	localPath := "./source.json"

	// Use wget to download the file
	cmd := exec.Command("wget", "-O", localPath, url)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to download file: %v\n", err)
	}
	sourceConfig := GetSourceConfig(localPath)

	return sourceConfig, nil
}

// GetSourceList -
func (s *Client) GetSourceList() ([]*SourceConfig, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("GetSourceList cost time: ", time.Since(start))
	}()

	localPath := "./source.json"
	sourceConfig := GetSourceConfig(localPath)

	return sourceConfig, nil
}

// StartTransmission -
func (s *Client) StartTransmission(param *Req) error {

	s.logger.Info("start tramsmission param: ", param)
	sourceDir := s.config.Dir
	remoteTarget, err := NewRemoteTarget(param, s.logger)
	if err != nil {
		s.logger.Errorf("Failed to create SSH client at target IP: %v", err)
		return err
	}

	defer remoteTarget.Close()

	err = remoteTarget.CheckRemoteDirectory(sourceDir)
	if err != nil {
		remoteTarget.redirect.Printf("Failed to check remote directory: %v", err)
		s.logger.Info("err: ", err)
		return err
	}

	/*--------------------get source path------------------------*/

	sourceConfig := GetSourceConfigByProjectName(param.ProjectName)
	if sourceConfig == nil {
		s.logger.Errorf("Failed to get source config by project name: %s", param.ProjectName)
		return fmt.Errorf("failed to get source config by project name: %s", param.ProjectName)
	}

	s.logger.Info("source: ", sourceConfig)

	remoteTargetSourcePath := path.Join(sourceDir, path.Base(sourceConfig.Source))

	remoteTarget.redirect.Printf("\remoteTargetSourcePath: %s\n", remoteTargetSourcePath)
	s.logger.Infof("remoteTargetSourcePath: %s\n", remoteTargetSourcePath)

	remoteSourceFile := sourceConfig.Source

	/*--------------------download file------------------------*/

	err = remoteTarget.DownloadTarball2Remote(remoteSourceFile, sourceDir)
	if err != nil {
		remoteTarget.redirect.Printf("Failed to download %s tarball to remote: %v", remoteSourceFile, err)
		s.logger.Info("err: ", err)
		return err
	}

	remoteTarget.redirect.Printf("Package transfer executed successfully.")
	s.logger.Info("Package transfer executed successfully.")

	/*---------------------extract tarball and backup the preversion file ----------------------*/

	backUpDir := path.Join(sourceDir, "old")
	// source tarbal
	err = remoteTarget.TransmitRemoteSource(remoteTargetSourcePath, sourceDir, backUpDir)
	if err != nil {
		remoteTarget.redirect.Printf("Failed to transmit source: %v", err)
		s.logger.Info("err: ", err)
		return err
	}

	/*---------------------import mysql script-----------------------*/

	// mysqlImage := "mysql:5.7"
	// mysqlContainerID, err := remoteTarget.CheckRemoteRunningContainer(mysqlImage)
	// if err != nil {
	// 	fs.logger.Info("Error checking MySQL container:", err)
	// 	return err
	// }
	// s.logger.Info("MySQL container ID:", mysqlContainerID)
	// remoteTarget.redirect.Printf("MySQL container ID: %s\n", mysqlContainerID)

	// if mysqlContainerID != "" {
	// 	remoteTarget.redirect.Printf("MySQL container is running. Importing sqlScript to database.")
	// 	s.logger.Info("MySQL container is running. Importing sqlScript to database.")
	// 	modelConfigItems := GetModelConfig()
	// 	err = remoteTarget.ImportSQLToRemoteDatabase(mysqlContainerID, remoteSQLScript, modelConfigItems)
	// 	if err != nil {
	// 		remoteTarget.redirect.Printf("Error importing sqlScript to MySQL: %v\n", err)
	// 		s.logger.Info("Error importing sqlScript to MySQL:", err)
	// 		return err
	// 	}
	// }

	/*---------------------remove tarball-----------------------*/

	// source tarbal
	err = remoteTarget.RemoveRemoteFile(remoteTargetSourcePath)
	if err != nil {
		remoteTarget.redirect.Printf("Failed to remove tarball: %v", err)
		s.logger.Info("err: ", err)
		return err
	}

	/*----------------------start container----------------------*/

	// // start the container
	// err = remoteTarget.StartRemoteContainer(*sourceConfig)
	// if err != nil {
	// 	remoteTarget.redirect.Printf("Failed to start container: %v", err)
	// 	s.logger.Info("err: ", err)
	// 	return err
	// }

	/*--------------------------------------------*/
	s.logger.Info("Transmission completed successfully")
	remoteTarget.redirect.Printf("Transmission completed successfully")

	return nil
}

// Detail -
type Detail struct {
	Log    string        `yaml:"log"`
	Config *SourceConfig `yaml:"config"`
	Date   string        `yaml:"date"`
}

// GetLogDetail -
func (s *Client) GetLogDetail(projectName, date string) (*Detail, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("GetTransmissionDetail cost time: ", time.Since(start))
	}()

	baseDir := "./upgrade_log"
	projectDirs, err := utils.ListDir(baseDir)
	if err != nil {
		return nil, err
	}

	var logPath string
	for _, projectDir := range projectDirs {
		if projectDir == projectName {
			dateDirs, err := utils.ListDir(filepath.Join(baseDir, projectDir))
			if err != nil {
				return nil, err
			}
			for _, dateDir := range dateDirs {
				if dateDir == date {
					logPath = filepath.Join(baseDir, projectDir, dateDir, "upgrade.log")
					break
				}
			}
			break
		}
	}

	if logPath == "" {
		return nil, fmt.Errorf("log not found for project: %s, date: %s", projectName, date)
	}

	log, err := os.ReadFile(logPath)
	if err != nil {
		return nil, err
	}

	sourceConfig := GetSourceConfigByProjectName(projectName)

	detail := &Detail{
		Log:    string(log),
		Config: sourceConfig,
		Date:   date,
	}

	return detail, nil
}

// LogList -
type LogList struct {
	ProjectName string `yaml:"project_name"`
	Date        string `yaml:"date"`
}

// GetLogList -
func (s *Client) GetLogList() ([]*LogList, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("GetUpgradeList cost time: ", time.Since(start))
	}()

	dir := "./upgrade_log"
	sourceDirs, err := utils.ListDir(dir)
	if err != nil {
		return nil, err
	}

	list := make([]*LogList, 0)
	for _, sourceDir := range sourceDirs {
		sourcePath := filepath.Join(dir, sourceDir)
		dateDirs, err := utils.ListDir(sourcePath)
		if err != nil {
			return nil, err
		}
		for _, dateDir := range dateDirs {
			list = append(list, &LogList{
				ProjectName: sourceDir,
				Date:        dateDir,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		dateI, _ := time.Parse("20060102150405", list[i].Date)
		dateJ, _ := time.Parse("20060102150405", list[j].Date)
		return dateI.After(dateJ)
	})

	return list, nil
}

// DeleteLogDetail -
func (s *Client) DeleteLogDetail(projectName, date string) error {
	projectDir := filepath.Join("./upgrade_log", projectName)
	if utils.IsFolderNotExist(projectDir) {
		return fmt.Errorf("project folder does not exist: %s", projectName)
	}

	dateDirs, err := utils.ListDir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to list directories in project folder: %s, error: %w", projectName, err)
	}

	for _, dateDir := range dateDirs {
		if dateDir == date {
			dir := filepath.Join(projectDir, dateDir)
			if err := utils.DeleteFolder(dir); err != nil {
				return fmt.Errorf("failed to delete upgrade log: %s, error: %w", date, err)
			}
			return nil
		}
	}

	return fmt.Errorf("upgrade log not found for date: %s", date)
}
