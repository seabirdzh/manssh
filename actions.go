package main

import (
	"fmt"
	"io/ioutil"

	"github.com/mikkeloscar/sshconfig"
	"github.com/urfave/cli"
)

var (
	path string
)

func list(c *cli.Context) error {
	hosts, _ := sshconfig.ParseSSHConfig(path)
	if len(c.Args()) > 0 {
		searchHosts := []*sshconfig.SSHHost{}
		for _, host := range hosts {
			values := []string{host.HostName, host.User, fmt.Sprintf("%d", host.Port)}
			values = append(values, host.Host...)
			if query(values, c.Args()) {
				searchHosts = append(searchHosts, host)
			}
		}
		hosts = searchHosts
	}
	printSuccessFlag()
	whiteBoldColor.Printf("Display %d records.\n\n", len(hosts))
	for _, host := range hosts {
		printHost(host)
	}
	return nil
}

func add(c *cli.Context) error {
	if err := argumentsCheck(c, 2, 2); err != nil {
		return err
	}
	newAlias := c.Args().Get(0)
	hostStr := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	if _, ok := hostMap[newAlias]; ok {
		printErrorFlag()
		return cli.NewExitError(fmt.Sprintf("'%s' ssh alias already exists", newAlias), 1)
	}
	host := parseHost(newAlias, hostStr)
	hosts = append(hosts, host)
	if err := saveHosts(hosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("'%s' alias config added successfully\n\n", newAlias)
	printHost(host)
	return nil
}

func update(c *cli.Context) error {
	fmt.Println("update command")
	return nil
}

func delete(c *cli.Context) error {
	if err := argumentsCheck(c, 1, -1); err != nil {
		return err
	}
	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	for _, alias := range c.Args() {
		if _, ok := hostMap[alias]; !ok {
			printErrorFlag()
			return cli.NewExitError(fmt.Sprintf("'%s' ssh alias not found", alias), 1)
		}
	}
	newHosts := []*sshconfig.SSHHost{}
	for _, host := range hosts {
		newAlias := []string{}
		for _, hostAlias := range host.Host {
			isDelete := false
			for _, deleteAlias := range c.Args() {
				if hostAlias == deleteAlias {
					isDelete = true
					break
				}
			}
			if !isDelete {
				newAlias = append(newAlias, hostAlias)
			}
		}
		host.Host = newAlias
		if len(host.Host) > 0 {
			newHosts = append(newHosts, host)
		}
	}
	if err := saveHosts(newHosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("deleted '%d' alias config\n", len(c.Args()))
	return nil
}

func rename(c *cli.Context) error {
	if err := argumentsCheck(c, 2, 2); err != nil {
		return err
	}

	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	oldName := c.Args().Get(0)
	newName := c.Args().Get(1)
	if _, ok := hostMap[oldName]; !ok {
		printErrorFlag()
		return cli.NewExitError("old ssh alias not found", 1)
	}
	if _, ok := hostMap[newName]; ok {
		printErrorFlag()
		return cli.NewExitError("new ssh alias already exists", 1)
	}
	host := hostMap[oldName]
	for i, name := range host.Host {
		if name == oldName {
			host.Host[i] = newName
			break
		}
	}
	if err := saveHosts(hosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("Rename from '%s' to '%s'\n\n", oldName, newName)
	printHost(host)
	return nil
}

func backup(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 1); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	backupPath := c.Args().First()
	err = ioutil.WriteFile(backupPath, data, 0644)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	whiteBoldColor.Printf("backup ssh config to '%s' success", backupPath)
	return nil
}