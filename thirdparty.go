// Copyright vinegar-development 2023

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	DXVKURL        = "https://github.com/doitsujin/dxvk/releases/download/v2.1/dxvk-2.1.tar.gz"
	RCOFFLAGSURL   = "https://raw.githubusercontent.com/L8X/Roblox-Client-Optimizer/main/ClientAppSettings.json"
	FPSUNLOCKERURL = "https://github.com/axstin/rbxfpsunlocker/releases/download/v4.4.4/rbxfpsunlocker-x64.zip"
)

func DxvkInstall() {
	dxvkTarballPath := filepath.Join(Dirs.Cache, "dxvk-2.0.tar.gz")

	Download(DXVKURL, dxvkTarballPath)

	dxvkTarball, err := os.Open(dxvkTarballPath)
	Errc(err)
	defer dxvkTarball.Close()

	dxvkGzip, err := gzip.NewReader(dxvkTarball)
	Errc(err)
	defer dxvkGzip.Close()

	dxvkTar := tar.NewReader(dxvkGzip)

	var dirInstall string
	for {
		header, err := dxvkTar.Next()

		if err == io.EOF {
			break
		}

		Errc(err)

		dllFile := path.Base(header.Name)
		archDir := path.Base(path.Dir(header.Name))

		switch header.Typeflag {
		case tar.TypeReg:
			if archDir == "x32" {
				dirInstall = "syswow64"
			} else {
				dirInstall = "system32"
			}

			writer, err := os.Create(filepath.Join(Dirs.Pfx, "drive_c", "windows", dirInstall, dllFile))
			Errc(err)
			log.Println("Gzipped:", writer.Name())
			io.Copy(writer, dxvkTar)
			writer.Close()
		}
	}

	Errc(os.RemoveAll(dxvkTarballPath))
}

func DxvkUninstall() {
	for _, dir := range []string{"syswow64", "system32"} {
		for _, dll := range []string{"d3d9", "d3d10core", "d3d11", "dxgi"} {
			dllFile := filepath.Join(Dirs.Pfx, "drive_c", "windows", dir, dll+".dll")
			log.Println("Removing DLL:", dllFile)
			Errc(os.RemoveAll(dllFile))
		}
	}
}

// Launch or automatically install axstin's rbxfpsunlocker.
// This function will also create it's own settings for rbxfpsunlocker, for
// faster or cleaner startup.
func RbxFpsUnlocker() {
	fpsUnlockerPath := filepath.Join(Dirs.Data, "rbxfpsunlocker.exe")
	_, err := os.Stat(fpsUnlockerPath)

	log.Println(err)
	if os.IsNotExist(err) {
		fpsUnlockerZip := filepath.Join(Dirs.Cache, "rbxfpsunlocker.zip")
		log.Println("Installing rbxfpsunlocker")
		Download(FPSUNLOCKERURL, fpsUnlockerZip)
		Unzip(fpsUnlockerZip, fpsUnlockerPath)
	}

	var settings = []string{
		"UnlockClient=true",
		"UnlockStudio=true",
		"FPSCapValues=[30.000000, 60.000000, 75.000000, 120.000000, 144.000000, 165.000000, 240.000000, 360.000000]",
		"FPSCapSelection=0",
		"FPSCap=0.000000",
		"CheckForUpdates=false",
		"NonBlockingErrors=true",
		"SilentErrors=true",
		"QuickStart=true",
	}

	settingsFile, err := os.Create(filepath.Join(Dirs.Cache, "settings"))
	Errc(err)
	defer settingsFile.Close()

	// FIXME: compare settings file, to check if user has modified the settings file
	log.Println("Writing custom rbxfpsunlocker settings to", settingsFile.Name())
	for _, setting := range settings {
		_, err := fmt.Fprintln(settingsFile, setting+"\r")
		Errc(err)
	}

	log.Println("Launching FPS Unlocker")
	Exec("wine", true, fpsUnlockerPath)

	// Since this file is always overwritten, just remove it anyway.
	Errc(os.RemoveAll(settingsFile.Name()))
}

// Download RCO (Roblox-Client-Optimizer)'s FFlags to the FFlags file provided.
func ApplyRCOFFlags(file string) {
	log.Println("Applying RCO FFlags")
	Download(RCOFFLAGSURL, file)
}
