package updater

import (
    "archive/tar"
    "archive/zip"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

const (
    githubAPIURL  = "https://api.github.com/repos/OlaHulleberg/codzure/releases/latest"
    githubRepoURL = "https://github.com/OlaHulleberg/codzure"
)

type GitHubRelease struct {
    TagName string `json:"tag_name"`
    Assets  []struct {
        Name               string `json:"name"`
        BrowserDownloadURL string `json:"browser_download_url"`
    } `json:"assets"`
}

func CheckForUpdates(currentVersion string) {
    if currentVersion == "dev" { return }
    latestVersion, err := getLatestVersion()
    if err != nil { return }
    if latestVersion != currentVersion && latestVersion != "" {
        fmt.Fprintf(os.Stderr, "\n⚠️  New version available: %s (current: %s)\n", latestVersion, currentVersion)
        fmt.Fprintf(os.Stderr, "   Run 'codzure update' to upgrade\n\n")
    }
}

func Update(currentVersion string) error {
    if currentVersion == "dev" { return fmt.Errorf("cannot update development build") }
    fmt.Println("Checking for updates...")
    release, err := getLatestRelease(); if err != nil { return fmt.Errorf("failed to check for updates: %w", err) }
    latest := release.TagName
    if latest == currentVersion { fmt.Printf("Already on latest version: %s\n", currentVersion); return nil }
    fmt.Printf("New version available: %s (current: %s)\n", latest, currentVersion)
    assetName := getBinaryAssetName()
    var url string
    for _, a := range release.Assets { if a.Name == assetName { url = a.BrowserDownloadURL; break } }
    if url == "" { return fmt.Errorf("no binary found for platform %s/%s", runtime.GOOS, runtime.GOARCH) }
    fmt.Printf("Downloading %s...\n", assetName)
    if err := downloadAndReplace(url); err != nil { return fmt.Errorf("failed to update: %w", err) }
    fmt.Printf("Successfully updated to version %s\n", latest)
    return nil
}

func getLatestVersion() (string, error) { r, e := getLatestRelease(); if e != nil { return "", e }; return r.TagName, nil }

func getLatestRelease() (*GitHubRelease, error) {
    resp, err := http.Get(githubAPIURL); if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode) }
    var release GitHubRelease
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil { return nil, err }
    return &release, nil
}

func getBinaryAssetName() string {
    osname := runtime.GOOS
    arch := runtime.GOARCH
    name := fmt.Sprintf("codzure_%s_%s", osname, arch)
    if osname == "windows" { name += ".zip" } else { name += ".tar.gz" }
    return name
}

func downloadAndReplace(url string) error {
    resp, err := http.Get(url); if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK { return fmt.Errorf("download failed with status %d", resp.StatusCode) }
    tmpFile, err := os.CreateTemp("", "codzure-archive-*"); if err != nil { return err }
    tmpPath := tmpFile.Name(); defer os.Remove(tmpPath)
    if _, err := io.Copy(tmpFile, resp.Body); err != nil { tmpFile.Close(); return err }
    tmpFile.Close()
    var binPath string
    if strings.HasSuffix(url, ".zip") { binPath, err = extractFromZip(tmpPath) } else if strings.HasSuffix(url, ".tar.gz") { binPath, err = extractFromTarGz(tmpPath) } else { return fmt.Errorf("unsupported archive format") }
    if err != nil { return fmt.Errorf("failed to extract binary: %w", err) }
    defer os.Remove(binPath)
    if err := os.Chmod(binPath, 0755); err != nil { return err }
    current, err := os.Executable(); if err != nil { return err }
    if runtime.GOOS == "windows" {
        backup := current + ".old"
        if err := os.Rename(current, backup); err != nil { return err }
        if err := os.Rename(binPath, current); err != nil { os.Rename(backup, current); return err }
        os.Remove(backup)
    } else {
        if err := os.Rename(binPath, current); err != nil { return err }
    }
    return nil
}

func extractFromTarGz(path string) (string, error) {
    f, err := os.Open(path); if err != nil { return "", err }
    defer f.Close()
    gz, err := gzip.NewReader(f); if err != nil { return "", err }
    defer gz.Close()
    tr := tar.NewReader(gz)
    for {
        h, err := tr.Next(); if err == io.EOF { break } ; if err != nil { return "", err }
        if h.Typeflag == tar.TypeReg && filepath.Base(h.Name) == "codzure" {
            tmp, err := os.CreateTemp("", "codzure-binary-*"); if err != nil { return "", err }
            p := tmp.Name()
            if _, err := io.Copy(tmp, tr); err != nil { tmp.Close(); os.Remove(p); return "", err }
            tmp.Close(); return p, nil
        }
    }
    return "", fmt.Errorf("binary not found in archive")
}

func extractFromZip(path string) (string, error) {
    zr, err := zip.OpenReader(path); if err != nil { return "", err }
    defer zr.Close()
    for _, file := range zr.File {
        if filepath.Base(file.Name) == "codzure.exe" {
            rc, err := file.Open(); if err != nil { return "", err }
            defer rc.Close()
            tmp, err := os.CreateTemp("", "codzure-binary-*.exe"); if err != nil { return "", err }
            p := tmp.Name()
            if _, err := io.Copy(tmp, rc); err != nil { tmp.Close(); os.Remove(p); return "", err }
            tmp.Close(); return p, nil
        }
    }
    return "", fmt.Errorf("binary not found in archive")
}
