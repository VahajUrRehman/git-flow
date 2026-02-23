import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import * as child_process from 'child_process';
import { promisify } from 'util';

const exec = promisify(child_process.exec);

export class GitFlowTerminal {
    private terminal: vscode.Terminal | undefined;
    private binaryPath: string | undefined;

    constructor() {
        this.binaryPath = this.findBinary();
    }

    /**
     * Find the gitflow-tui binary
     */
    private findBinary(): string | undefined {
        const config = vscode.workspace.getConfiguration('gitflow');
        const configuredPath = config.get<string>('binaryPath');

        if (configuredPath && fs.existsSync(configuredPath)) {
            return configuredPath;
        }

        // Try to find in PATH
        const possibleNames = ['gitflow-tui', 'gitflow'];
        
        for (const name of possibleNames) {
            try {
                const result = child_process.execSync(`which ${name}`, { encoding: 'utf-8' });
                const binaryPath = result.trim();
                if (binaryPath && fs.existsSync(binaryPath)) {
                    return binaryPath;
                }
            } catch {
                // Not found, continue
            }
        }

        // Try common installation paths
        const commonPaths = [
            path.join(process.env.HOME || '', '.local', 'bin', 'gitflow-tui'),
            path.join(process.env.HOME || '', 'bin', 'gitflow-tui'),
            '/usr/local/bin/gitflow-tui',
            '/usr/bin/gitflow-tui',
            path.join(process.env.ProgramFiles || '', 'gitflow-tui', 'gitflow-tui.exe'),
        ];

        for (const p of commonPaths) {
            if (fs.existsSync(p)) {
                return p;
            }
        }

        return undefined;
    }

    /**
     * Get the workspace folder
     */
    private getWorkspaceFolder(): string | undefined {
        const folders = vscode.workspace.workspaceFolders;
        if (folders && folders.length > 0) {
            return folders[0].uri.fsPath;
        }
        return undefined;
    }

    /**
     * Check if we're in a git repository
     */
    private async isGitRepository(folder: string): Promise<boolean> {
        try {
            await exec('git rev-parse --git-dir', { cwd: folder });
            return true;
        } catch {
            return false;
        }
    }

    /**
     * Open GitFlow TUI in terminal
     */
    public async open(): Promise<void> {
        if (!this.binaryPath) {
            const result = await vscode.window.showErrorMessage(
                'GitFlow TUI binary not found. Please install it or configure the path.',
                'Configure',
                'Install'
            );

            if (result === 'Configure') {
                vscode.commands.executeCommand('workbench.action.openSettings', 'gitflow.binaryPath');
            } else if (result === 'Install') {
                vscode.env.openExternal(vscode.Uri.parse('https://github.com/gitflow/tui#installation'));
            }
            return;
        }

        const workspaceFolder = this.getWorkspaceFolder();
        if (!workspaceFolder) {
            vscode.window.showErrorMessage('No workspace folder open');
            return;
        }

        if (!await this.isGitRepository(workspaceFolder)) {
            const result = await vscode.window.showWarningMessage(
                'Current folder is not a git repository. Initialize?',
                'Yes',
                'No'
            );

            if (result === 'Yes') {
                try {
                    await exec('git init', { cwd: workspaceFolder });
                    vscode.window.showInformationMessage('Git repository initialized');
                } catch (error) {
                    vscode.window.showErrorMessage(`Failed to initialize git: ${error}`);
                    return;
                }
            } else {
                return;
            }
        }

        // Get configuration
        const config = vscode.workspace.getConfiguration('gitflow');
        const useIntegrated = config.get<boolean>('terminal.integrated', true);

        if (useIntegrated) {
            this.openInIntegratedTerminal(workspaceFolder);
        } else {
            this.openInExternalTerminal(workspaceFolder);
        }
    }

    /**
     * Open in integrated terminal
     */
    private openInIntegratedTerminal(cwd: string): void {
        // Close existing terminal
        if (this.terminal) {
            this.terminal.dispose();
        }

        // Create new terminal
        this.terminal = vscode.window.createTerminal({
            name: 'GitFlow TUI',
            cwd: cwd,
            shellPath: this.binaryPath,
            shellArgs: [],
            env: {
                ...process.env,
                GITFLOW_VSCODE: '1',
            },
        });

        this.terminal.show(true);

        // Handle terminal close
        vscode.window.onDidCloseTerminal((t) => {
            if (t === this.terminal) {
                this.terminal = undefined;
            }
        });
    }

    /**
     * Open in external terminal
     */
    private openInExternalTerminal(cwd: string): void {
        const config = vscode.workspace.getConfiguration('gitflow');
        
        // Build command with environment variables
        const env = {
            ...process.env,
            GITFLOW_VSCODE: '1',
        };

        const command = `cd "${cwd}" && ${this.binaryPath}`;

        // Open external terminal based on platform
        if (process.platform === 'win32') {
            child_process.exec(`start cmd /k "${command}"`, { env });
        } else if (process.platform === 'darwin') {
            child_process.exec(`osascript -e 'tell application "Terminal" to do script "${command}"' -e 'tell application "Terminal" to activate'`);
        } else {
            // Linux - try common terminal emulators
            const terminals = [
                'gnome-terminal',
                'konsole',
                'xfce4-terminal',
                'xterm',
                'alacritty',
                'kitty',
            ];

            for (const term of terminals) {
                try {
                    child_process.execSync(`which ${term}`);
                    let termCommand = '';
                    
                    switch (term) {
                        case 'gnome-terminal':
                        case 'xfce4-terminal':
                            termCommand = `${term} --working-directory="${cwd}" -- ${this.binaryPath}`;
                            break;
                        case 'konsole':
                            termCommand = `${term} --workdir "${cwd}" -e ${this.binaryPath}`;
                            break;
                        case 'alacritty':
                        case 'kitty':
                            termCommand = `${term} --working-directory="${cwd}" -e ${this.binaryPath}`;
                            break;
                        default:
                            termCommand = `cd "${cwd}" && ${term} -e ${this.binaryPath}`;
                    }

                    child_process.exec(termCommand, { env });
                    break;
                } catch {
                    continue;
                }
            }
        }
    }

    /**
     * Toggle terminal visibility
     */
    public toggle(): void {
        if (this.terminal) {
            this.close();
        } else {
            this.open();
        }
    }

    /**
     * Close the terminal
     */
    public close(): void {
        if (this.terminal) {
            this.terminal.dispose();
            this.terminal = undefined;
        }
    }

    /**
     * Run a git command
     */
    public async runCommand(command: string): Promise<void> {
        const workspaceFolder = this.getWorkspaceFolder();
        if (!workspaceFolder) {
            vscode.window.showErrorMessage('No workspace folder open');
            return;
        }

        try {
            const { stdout, stderr } = await exec(`git ${command}`, { cwd: workspaceFolder });
            
            if (stderr) {
                vscode.window.showWarningMessage(stderr);
            }

            if (stdout) {
                // Show output in output channel
                const outputChannel = vscode.window.createOutputChannel('GitFlow');
                outputChannel.appendLine(stdout);
                outputChannel.show();
            }

            // Refresh views
            vscode.commands.executeCommand('gitflow.refresh');
        } catch (error) {
            vscode.window.showErrorMessage(`Git command failed: ${error}`);
        }
    }

    /**
     * Dispose resources
     */
    public dispose(): void {
        this.close();
    }
}
