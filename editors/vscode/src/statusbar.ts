import * as vscode from 'vscode';
import * as child_process from 'child_process';

export class GitFlowStatusBar {
    private statusBarItem: vscode.StatusBarItem;

    constructor() {
        this.statusBarItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            100
        );
        this.statusBarItem.command = 'gitflow.toggle';
        this.update();
    }

    public register(context: vscode.ExtensionContext): void {
        context.subscriptions.push(this.statusBarItem);
        
        // Update on file save
        vscode.workspace.onDidSaveTextDocument(() => {
            this.update();
        }, null, context.subscriptions);

        // Update on git events
        const gitExtension = vscode.extensions.getExtension('vscode.git');
        if (gitExtension) {
            gitExtension.activate().then(git => {
                const api = git.exports.getAPI(1);
                api.onDidChangeState(() => {
                    this.update();
                });
            });
        }

        // Update every 10 seconds
        setInterval(() => {
            this.update();
        }, 10000);

        this.statusBarItem.show();
    }

    public update(): void {
        try {
            const workspaceFolder = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
            if (!workspaceFolder) {
                this.statusBarItem.hide();
                return;
            }

            // Get current branch
            let branch = '';
            try {
                branch = child_process.execSync('git branch --show-current', {
                    cwd: workspaceFolder,
                    encoding: 'utf-8'
                }).trim();
            } catch {
                this.statusBarItem.hide();
                return;
            }

            // Get status
            let status = '';
            try {
                const result = child_process.execSync('git status --porcelain', {
                    cwd: workspaceFolder,
                    encoding: 'utf-8'
                }).trim();
                
                if (result) {
                    const lines = result.split('\n').filter(line => line.trim());
                    const staged = lines.filter(l => l[0] !== ' ' && l[0] !== '?').length;
                    const unstaged = lines.filter(l => l[1] !== ' ').length;
                    const untracked = lines.filter(l => l[0] === '?').length;

                    const parts: string[] = [];
                    if (staged > 0) parts.push(`+${staged}`);
                    if (unstaged > 0) parts.push(`~${unstaged}`);
                    if (untracked > 0) parts.push(`?${untracked}`);
                    
                    if (parts.length > 0) {
                        status = ` ${parts.join(' ')}`;
                    }
                }
            } catch {
                // Ignore
            }

            // Get ahead/behind
            let sync = '';
            try {
                const result = child_process.execSync('git rev-list --left-right --count HEAD...@{upstream}', {
                    cwd: workspaceFolder,
                    encoding: 'utf-8'
                }).trim();
                
                const [ahead, behind] = result.split('\t').map(n => parseInt(n, 10));
                if (ahead > 0 || behind > 0) {
                    sync = ` ↑${ahead}↓${behind}`;
                }
            } catch {
                // No upstream
            }

            // Update status bar
            const config = vscode.workspace.getConfiguration('gitflow');
            const primaryColor = config.get<string>('theme.primary', '#00D9A5');
            
            this.statusBarItem.text = `$(git-branch) ${branch}${status}${sync}`;
            this.statusBarItem.tooltip = `GitFlow: ${branch}\nClick to open GitFlow TUI`;
            
            if (status) {
                this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
            } else {
                this.statusBarItem.backgroundColor = undefined;
            }

            this.statusBarItem.show();
        } catch (error) {
            console.error('Error updating status bar:', error);
            this.statusBarItem.hide();
        }
    }

    public dispose(): void {
        this.statusBarItem.dispose();
    }
}
