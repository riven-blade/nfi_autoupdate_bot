package src

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
	"os"
	"strings"
)

// GitClone 运行 git clone 命令，使用 go-git 库
func GitClone(repoURL, targetDir string) error {
	_, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	fmt.Println("Repository cloned successfully!")
	return nil
}

// GitPull 运行 git pull 命令
func GitPull(repoDir string) error {
	// 打开本地仓库
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	// 获取工作树 (worktree)
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	// 拉取最新的代码
	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",  // 指定远程仓库名
		Progress:   os.Stdout, // 显示拉取进度
		Force:      true,      // 如果有冲突，可以强制拉取
	})

	// 如果没有新的更新，会返回 `already up-to-date` 错误
	if errors.Is(git.NoErrAlreadyUpToDate, err) {
		fmt.Println("Repository is already up-to-date")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to pull repository: %v", err)
	}

	fmt.Println("Repository pulled successfully!")
	return nil
}

// CalculateCheckSum / 计算文件的 SHA256 校验和
func CalculateCheckSum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CheckAndReplaceFile 检查和替换文件
func CheckAndReplaceFile(localPath, githubFileURL string, isChange bool) (error, bool) {
	// 计算本地文件的校验和
	localCheckSum, err := CalculateCheckSum(localPath)
	if err != nil {
		return fmt.Errorf("error calculating local file checksum: %v", err), false
	}

	// 计算下载文件的校验和
	githubCheckSum, err := CalculateCheckSum(githubFileURL)
	if err != nil {
		return fmt.Errorf("error calculating GitHub file checksum: %v", err), false
	}

	fmt.Printf("Local CheckSum: %s\n", localCheckSum)
	fmt.Printf("GitHub CheckSum: %s\n", githubCheckSum)

	// 比较校验和
	if localCheckSum == githubCheckSum {
		fmt.Println("Files are identical. No update needed.")
		return nil, false
	} else if isChange {
		// 如果不同，覆盖本地文件
		err := copyFile(githubFileURL, localPath)
		if err != nil {
			return fmt.Errorf("failed to replace local file: %v", err), false
		}
		fmt.Println("Local file replaced with the latest version from GitHub.")
	}

	return nil, true
}

func CheckFolderExists(folderPath string) (bool, error) {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// copyFile 简化文件复制函数
func copyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 确保所有数据都写入磁盘
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func GetVersions(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	targetLine := 69
	var targetString string

	for scanner.Scan() {
		lineNumber++
		if lineNumber == targetLine {
			line := scanner.Text()
			if strings.Contains(line, `return "v`) {
				start := strings.Index(line, `"`)
				end := strings.LastIndex(line, `"`)
				if start != -1 && end != -1 && start != end {
					targetString = line[start+1 : end]
				}
			}
			break
		}
	}
	return targetString
}
