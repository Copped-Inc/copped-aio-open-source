package update

import (
	"fmt"
	"github.com/Copped-Inc/aio-types/console"
	"os/exec"
)

func Update() error {

	console.Log("Update Runner Image")
	// INSERT GITHUBHACCESSTOKEN here and correct COMPANY
	curl := fmt.Sprintf("curl \"https://x-access-token:GITHUBHACCESSTOKEN@raw.githubusercontent.com/COMPANY/client-local-v3/main/target/release/client-local-v3\" --output /home/runner/aio-instances-api/_work/aio-instances-api/aio-instances-api/client/client-local-v3")
	prune := fmt.Sprintf("docker system prune -a -f")
	del := fmt.Sprintf("docker image rm client -f")
	build := fmt.Sprintf("docker build -t client /home/runner/aio-instances-api/_work/aio-instances-api/aio-instances-api/client")

	_, err := exec.Command("/bin/sh", "-c", curl).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/sh", "-c", prune).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/sh", "-c", del).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/sh", "-c", build).Output()
	return err

}
