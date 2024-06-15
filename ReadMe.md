# Go Kube Watcher

Go version 1.22.4

---

This application uses users kubeconfig to connect and watch kubernetes pod changes and will write the change to the user file of "~/.kubeWatcher.result"

The "~/.kubeWatcher.result" gets overwritten every time the application starts

If running on Windows, the application can be installed and exected in task scheduler.

Steps: 

To build 

`

    go build .
    # This produces goKubeWatcher.exe

`
Place the output to a directory you manage.
Then open task scheduler with admin privileges
select your trigger 
select the build executable

you can then observe

"~/.kubeWatcher.result"
