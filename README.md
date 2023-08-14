# drain-q

Organic & environmentally friendly queue drainer

### Executing the binary rom anywhere

To make a Go binary executable callable from anywhere in the terminal, you need to place it in a directory that is included in your system's PATH environment variable. This way, when you type the binary's name in the terminal, the system will be able to locate and execute it. Here's what you can do:

Create a Directory for Your Binaries:
Create a directory where you can store your personal Go binaries. You can name it anything you like, such as bin, go_bin, or even your username, like ~/username/bin.

Add the Directory to PATH:
You need to add the directory where you placed your binaries to your system's PATH environment variable. This allows the system to find and execute binaries in that directory.

Temporary Update (for the current terminal session):
You can update the PATH temporarily for the current terminal session using the export command. Replace your_directory_path with the actual path of the directory where you stored your binaries:

```sh
export PATH=$PATH:your_directory_path
```
Permanent Update (for future sessions):
To make this change permanent, add the following line to your shell profile configuration file. The specific file depends on the shell you are using. For example, for the Bash shell:

```sh
echo 'export PATH=$PATH:your_directory_path' >> ~/.bashrc
```
After adding the line, you might need to restart your terminal or run source ~/.bashrc (or the appropriate command for your shell) to apply the changes.

Move Your Binary:
Move your compiled Go binary into the directory you've added to your PATH.

Make the Binary Executable:
Ensure that the binary has executable permissions. You can use the chmod command to make it executable:

```sh
mv drainq your_directory_path
```
Now, when you type the name of your binary in the terminal, regardless of your current working directory, your system will be able to locate and execute it because the directory is part of the PATH.

Remember that modifying system-wide PATH variables requires administrative privileges on some systems. Always exercise caution when making changes to environment variables.
