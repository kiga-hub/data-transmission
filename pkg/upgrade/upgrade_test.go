package upgrade

// // TestUpgrade -
// func TestUpgrade(t *testing.T) {
// 	deployDir := "/data/app/tritonserver/"

// 	user := "root"
// 	password := "123456"
// 	sourceIP := "10.1.3.31:22"
// 	sshSourceServer, err := NewUpgradeClient(user, password, sourceIP)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	defer sshSourceServer.client.Close()

// 	targetIP := "10.1.3.46:22"
// 	sshTargetClient, err := NewUpgradeClient(user, password, targetIP)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	defer sshTargetClient.client.Close()

// 	err = sshTargetClient.CheckRemoteDirectory(deployDir, sshTargetClient)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 	}

// 	/*--------------------------------------------*/

// 	sourceConfig := GetSourceConfig()

// 	targetModelLibPath := path.Join(deployDir, path.Base(sourceConfig.ModelLib))
// 	targetModelRepository := path.Join(deployDir, path.Base(sourceConfig.ModelRepository))
// 	targetSavedSWKWav := path.Join(deployDir, path.Base(sourceConfig.SavedSwkWav))
// 	targetSqlScript := path.Join(deployDir, path.Base(sourceConfig.SqlScript))

// 	fmt.Printf("targetModelLibPath: %s\ntargetModelRepository: %s\ntargetSavedSWKWav: %s\ntargetSqlScript: %s\n", targetModelLibPath, targetModelRepository, targetSavedSWKWav, targetSqlScript)

// 	remoteModelLibFile := sourceConfig.ModelLib
// 	// remoteModeRepository := sourceConfig.ModelRepository
// 	// remoteSavedSWDWav := sourceConfig.SavedSwkWav
// 	remoteSqlScript := sourceConfig.SqlScript

// 	/*--------------------------------------------*/

// 	err = sshSourceServer.TransferTarball2Remote(sourceIP, targetIP, remoteModelLibFile, targetModelLibPath, sshTargetClient)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}

// 	err = sshSourceServer.TransferTarball2Remote(sourceIP, targetIP, remoteSqlScript, targetSqlScript, sshTargetClient)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	fmt.Println("Package transfer executed successfully.")

// 	/*--------------------------------------------*/

// 	// mysqlImage := "harbor.aithu.com:80/middleware/mysql:5.7"
// 	// mysqlContainerID, err := checkRemoteRunningContainer(mysqlImage, sshTargetClient)
// 	// if err != nil {
// 	// 	fmt.Println("Error checking MySQL container:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("MySQL container ID:", mysqlContainerID)

// 	// if mysqlContainerID != "" {
// 	// 	fmt.Println("MySQL container is running. Importing sqlScript to database.")
// 	// 	modelConfigItems := GetModelConfig()
// 	// 	err = ImportSQLToRemoteDatabase(mysqlContainerID, targetSqlScript, modelConfigItems, sshTargetClient)
// 	// 	if err != nil {
// 	// 		fmt.Println("Error importing sqlScript to MySQL:", err)
// 	// 		return
// 	// 	}
// 	// }

// 	/*--------------------------------------------*/

// 	err = sshTargetClient.RemoveRemoteFile(targetModelLibPath, sshTargetClient)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	/*--------------------------------------------*/

// 	// start the container
// 	err = sshTargetClient.StartRemoteContainer(sourceConfig)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	/*--------------------------------------------*/
// }
