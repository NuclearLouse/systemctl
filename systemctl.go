package systemctl

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

/*
Example:
	ok, err := IsRootPermissions()
	if err != nil {
		log.Fatalln("check root permissions:", err)
	}
	if !ok {
		log.Fatal("user is not root")
	}

	command := "stop"
	service := "rater-backup"
	out, err := ServiceWithRootPermissions(command, service)
	if err != nil {
		log.Fatalln("exec command:", err)
	}
	command = "status"
	out, err := ServiceWithRootPermissions(command, service)
	if err != nil {
		log.Fatalln("exec command:", err)
	}
	fmt.Println(string(out))
*/

func IsRootPermissions() (bool, error) {
	//You can also add a check for "whoami" and use the USER_NAME environment variable
	var e *exec.ExitError
	out, err := exec.Command("id", "-u").Output()
	if err != nil && !errors.As(err, &e) {
		return false, fmt.Errorf("check root permissions: %w", err)
	}
	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(out[:len(out)-1]))
	if err != nil {
		return false, fmt.Errorf("converting the response from the check root permissions: %w", err)
	}
	return i == 0, nil
}

//Сommand expected: start, stop or status.
//After the "start" or "stop" command, it is desirable to call the "status" command
//Service file name can be passed without extension .service
func ServiceWithRootPermissions(command, service string) ([]byte, error) {
	var e *exec.ExitError
	out, err := exec.Command("systemctl", command, service).Output()
	if err != nil && !errors.As(err, &e) {
		return nil, err
	}
	return out, nil
}

func ServiceNoRootPermissions(command, service, password string) ([]byte, error) {
	//another version of the command "stop":
	// _, err := exec.Command("sh", "-c", "echo '"+ password +"' | sudo -S pkill -SIGINT "+service)

	var e *exec.ExitError
	out, err := exec.Command("echo", fmt.Sprintf("'%s'", password), "| sudo -S systemctl", command, service).Output()
	if err != nil && !errors.As(err, &e) {
		return nil, err
	}
	return out, nil
}
