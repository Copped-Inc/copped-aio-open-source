package update

import (
	"github.com/Copped-Inc/aio-types/console"
	"io"
	"net/http"
	"os"
)

func download(f file, r *http.Request) error {

	if r == nil {
		console.Log("Download", "File", f.name)
	} else {
		console.RequestLog(r, "Download", "File", f.name)
	}

	// INSERT GITHUBHACCESSTOKEN here and correct COMPANY
	req, err := http.NewRequest(http.MethodGet, "https://x-access-token:GITHUBHACCESSTOKEN@raw.githubusercontent.com/COMPANY"+f.path, nil)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	out, err := os.Create(f.name)
	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, res.Body)
	return err

}

func DownloadAll() error {
	err := download(ClientWindows, nil)
	if err != nil {
		return err
	}

	err = download(ClientLinux, nil)
	if err != nil {
		return err
	}

	err = download(ClientVersion, nil)
	if err != nil {
		return err
	}

	err = download(PaymentsWindows, nil)
	if err != nil {
		return err
	}

	err = download(PaymentsVersion, nil)
	return err
}

func DownloadClient(r *http.Request) error {
	err := download(ClientWindows, r)
	if err != nil {
		return err
	}

	err = download(ClientLinux, r)
	if err != nil {
		return err
	}

	err = download(ClientVersion, r)
	return err
}

func DownloadPayments(r *http.Request) error {
	err := download(PaymentsWindows, r)
	if err != nil {
		return err
	}

	err = download(PaymentsVersion, r)
	return err
}

type file struct {
	name string
	path string
}

type File file

var (
	ClientWindows file = file{
		name: "client-local-v3.exe",
		path: "/client-local-v3/main/target/debug/client-local-v3.exe",
	}
	ClientLinux file = file{
		name: "client-local-v3",
		path: "/client-local-v3/main/target/release/client-local-v3",
	}
	ClientVersion file = file{
		name: "version-client-local-v3",
		path: "/client-local-v3/main/version-client-local-v3",
	}
	PaymentsWindows file = file{
		name: "payments.exe",
		path: "/aio-payments/main/build/bin/payments.exe",
	}
	PaymentsVersion file = file{
		name: "version-payments",
		path: "/aio-payments/main/version-payments",
	}
)
