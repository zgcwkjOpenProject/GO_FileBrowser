package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/users"
	"golang.org/x/net/webdav"
)

// webdav 入口
func webdavUse(next http.Handler, prefix string, store *storage.Storage, server *settings.Server) http.Handler {
	// r.Use(func(next http.Handler) http.Handler {
	// 	return webdavUse(next, "/", store, server)
	// })
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 判断是否为浏览器UA
		userAgent := strings.ToLower(r.Header.Get("User-Agent"))
		isBrowser := strings.Contains(userAgent, "chrome") || strings.Contains(userAgent, "firefox") ||
			strings.Contains(userAgent, "safari") || strings.Contains(userAgent, "edge") ||
			strings.Contains(userAgent, "opera")

		// 继续处理浏览器访问
		if isBrowser {
			next.ServeHTTP(w, r)
			return
		}

		// 文件系统
		handle(webdavHandler, prefix, store, server)
	})
}

// webdav 文件系统
func webdavHandler(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// 使用程序内部用户认证
	user, err := verifyUserAuth(r, d)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		return http.StatusUnauthorized, nil
	}

	// 检查用户基本权限
	d.user = user
	userPath := path.Join(d.server.Root, d.user.Scope)
	if !d.user.Perm.Download {
		return http.StatusForbidden, nil
	}

	// 如果需要 GET/HEAD 目录时自动转 PROPFIND
	if r.Method == "GET" || r.Method == "HEAD" {
		relPath := strings.TrimPrefix(r.URL.Path, "/")
		absPath := filepath.Join(userPath, relPath)
		info, err := os.Stat(absPath)
		if err == nil && info.IsDir() {
			r.Method = "PROPFIND"
			if r.Header.Get("Depth") == "" {
				r.Header.Add("Depth", "1")
			}
		}
	}

	// WebDAV处理器
	handler := &webdav.Handler{
		FileSystem: &webdavFileSystem{
			fs:      webdav.Dir(userPath),
			checker: d.Check,
		},
		LockSystem: webdav.NewMemLS(),
	}

	// 处理请求
	handler.ServeHTTP(w, r)

	return 0, nil
}

// webdav 用户验证程序
func verifyUserAuth(r *http.Request, d *data) (*users.User, error) {
	// 获取帐号密码
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("no auth provided")
	}

	// 通过用户名查询用户
	user, err := d.store.Users.Get(d.server.Root, username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if !users.CheckPwd(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// webdavFileSystem 包装器，支持用户规则检查
type webdavFileSystem struct {
	fs      webdav.FileSystem
	checker func(string) bool
}

// 检查路径的每一级分量是否都符合规则
func (wfs *webdavFileSystem) checkPathAll(name string) bool {
	cleaned := filepath.Clean(name)
	if cleaned == "." || cleaned == "/" || cleaned == "" {
		return true
	}
	// 兼容绝对路径和相对路径
	cleaned = strings.TrimPrefix(cleaned, "/")
	parts := strings.Split(cleaned, string(filepath.Separator))
	cur := ""
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		cur = filepath.Join(cur, p)
		if !wfs.checker(cur) {
			return false
		}
	}
	return true
}

func (wfs *webdavFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if !wfs.checkPathAll(name) {
		return os.ErrPermission
	}
	return wfs.fs.Mkdir(ctx, name, perm)
}

func (wfs *webdavFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if !wfs.checkPathAll(name) {
		return nil, os.ErrPermission
	}
	return wfs.fs.OpenFile(ctx, name, flag, perm)
}

func (wfs *webdavFileSystem) RemoveAll(ctx context.Context, name string) error {
	if !wfs.checkPathAll(name) {
		return os.ErrPermission
	}
	return wfs.fs.RemoveAll(ctx, name)
}

func (wfs *webdavFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	if !wfs.checkPathAll(newName) {
		return os.ErrPermission
	}
	return wfs.fs.Rename(ctx, oldName, newName)
}

func (wfs *webdavFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if !wfs.checkPathAll(name) {
		return nil, os.ErrPermission
	}
	return wfs.fs.Stat(ctx, name)
}
