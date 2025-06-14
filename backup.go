package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"ms-admin/api/constants"
	loconfig "ms-admin/config"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func runBackups(db *mongo.Database) {
	mainDirBackup := os.Getenv("BACKUP_ROOT")
	if mainDirBackup == "" {
		mainDirBackup = "./backups"
	}

	go func() {
		dailyTicker := time.NewTicker(24 * time.Hour)
		for range dailyTicker.C {
			if err := Backup(db, "day", mainDirBackup); err != nil {
				log.Printf(constants.ErrBackUpDay, err)
			}
		}
	}()

	go func() {
		weeklyTicker := time.NewTicker(168 * time.Hour)
		for range weeklyTicker.C {
			if err := Backup(db, "week", mainDirBackup); err != nil {
				log.Printf(constants.ErrBackUpWeek, err)
			}
		}
	}()

	go func() {
		monthlyTicker := time.NewTicker(720 * time.Hour)
		for range monthlyTicker.C {
			if err := Backup(db, "month", mainDirBackup); err != nil {
				log.Printf(constants.ErrBackUpMonth, err)
			}
		}
	}()
}

func Backup(db *mongo.Database, backupType, mainDirBackup string) error {
	bkDir := filepath.Join(mainDirBackup, backupType)

	if err := os.MkdirAll(bkDir, 0755); err != nil {
		return fmt.Errorf(constants.ErrCreateFolderBackup, err)
	}

	backupFile := filepath.Join(
		bkDir,
		fmt.Sprintf("backup-%s-%s.tar.gz", backupType, time.Now().Format("2006-01-02-15-04-05")),
	)

	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf(constants.ErrCreateFileBackup, err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	collections, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf(constants.ErrCollectionsNotFound, err)
	}

	for _, collectionName := range collections {
		if err := exportCollection(db, collectionName, tarWriter); err != nil {
			return fmt.Errorf(constants.ErrExportCollection, collectionName, err)
		}
	}

	log.Printf(constants.SuccUploadBackup, backupFile)
	return nil
}

func exportCollection(db *mongo.Database, collectionName string, tarWriter *tar.Writer) error {
	collection := db.Collection(collectionName)

	cfg, err := loconfig.LoadLocalConfig()
	if err != nil {
		return err
	}

	cursor, err := collection.Find(context.Background(), bson.M{
		"status": bson.M{"$ne": constants.AdminStatus},
		"email":  bson.M{"$ne": cfg.ADMIN_EMAIL},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.json", collectionName))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	var docs []bson.M
	if err = cursor.All(context.Background(), &docs); err != nil {
		return err
	}

	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(docs); err != nil {
		return err
	}

	fileInfo, err := tmpFile.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    fmt.Sprintf("%s.json", collectionName),
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	}

	if err = tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		return err
	}

	if _, err = io.Copy(tarWriter, tmpFile); err != nil {
		return err
	}

	return nil
}
