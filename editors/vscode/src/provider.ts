import * as vscode from 'vscode';
import * as child_process from 'child_process';

export class GitFlowProvider implements vscode.TreeDataProvider<GitFlowItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<GitFlowItem | undefined | null | void> = new vscode.EventEmitter<GitFlowItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<GitFlowItem | undefined | null | void> = this._onDidChangeTreeData.event;

    private items: GitFlowItem[] = [];

    constructor() {
        this.refresh();
    }

    refresh(): void {
        this.loadData();
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: GitFlowItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: GitFlowItem): Thenable<GitFlowItem[]> {
        if (element) {
            return Promise.resolve(element.children || []);
        }
        return Promise.resolve(this.items);
    }

    private loadData(): void {
        this.items = [];

        // Add main actions
        this.items.push(
            new GitFlowItem(
                'Open GitFlow TUI',
                vscode.TreeItemCollapsibleState.None,
                {
                    command: 'gitflow.open',
                    title: 'Open GitFlow TUI',
                    arguments: []
                },
                '$(terminal)'
            )
        );

        // Add git status section
        const statusItem = new GitFlowItem(
            'Status',
            vscode.TreeItemCollapsibleState.Collapsed,
            undefined,
            '$(git-branch)'
        );
        statusItem.children = this.getStatusItems();
        this.items.push(statusItem);

        // Add recent commits
        const commitsItem = new GitFlowItem(
            'Recent Commits',
            vscode.TreeItemCollapsibleState.Collapsed,
            undefined,
            '$(git-commit)'
        );
        commitsItem.children = this.getCommitItems();
        this.items.push(commitsItem);

        // Add branches
        const branchesItem = new GitFlowItem(
            'Branches',
            vscode.TreeItemCollapsibleState.Collapsed,
            undefined,
            '$(git-branch)'
        );
        branchesItem.children = this.getBranchItems();
        this.items.push(branchesItem);

        // Add remotes
        const remotesItem = new GitFlowItem(
            'Remotes',
            vscode.TreeItemCollapsibleState.Collapsed,
            undefined,
            '$(repo)'
        );
        remotesItem.children = this.getRemoteItems();
        this.items.push(remotesItem);

        // Add stash
        const stashItem = new GitFlowItem(
            'Stash',
            vscode.TreeItemCollapsibleState.Collapsed,
            undefined,
            '$(archive)'
        );
        stashItem.children = this.getStashItems();
        this.items.push(stashItem);
    }

    private getStatusItems(): GitFlowItem[] {
        const items: GitFlowItem[] = [];
        
        try {
            const result = child_process.execSync('git status --porcelain', { encoding: 'utf-8' });
            const lines = result.split('\n').filter(line => line.trim());

            if (lines.length === 0) {
                items.push(new GitFlowItem(
                    'Working tree clean',
                    vscode.TreeItemCollapsibleState.None,
                    undefined,
                    '$(check)'
                ));
            } else {
                for (const line of lines.slice(0, 10)) {
                    const status = line.substring(0, 2);
                    const file = line.substring(3);
                    
                    let icon = '$(file)';
                    if (status.includes('M')) icon = '$(diff-modified)';
                    else if (status.includes('A')) icon = '$(diff-added)';
                    else if (status.includes('D')) icon = '$(diff-removed)';
                    else if (status.includes('?')) icon = '$(file-add)';

                    items.push(new GitFlowItem(
                        `${status} ${file}`,
                        vscode.TreeItemCollapsibleState.None,
                        undefined,
                        icon
                    ));
                }

                if (lines.length > 10) {
                    items.push(new GitFlowItem(
                        `... and ${lines.length - 10} more files`,
                        vscode.TreeItemCollapsibleState.None,
                        undefined,
                        '$(ellipsis)'
                    ));
                }
            }
        } catch {
            items.push(new GitFlowItem(
                'Not a git repository',
                vscode.TreeItemCollapsibleState.None,
                undefined,
                '$(error)'
            ));
        }

        return items;
    }

    private getCommitItems(): GitFlowItem[] {
        const items: GitFlowItem[] = [];
        
        try {
            const result = child_process.execSync('git log --oneline -10', { encoding: 'utf-8' });
            const lines = result.split('\n').filter(line => line.trim());

            for (const line of lines) {
                const hash = line.substring(0, 7);
                const message = line.substring(8);

                items.push(new GitFlowItem(
                    `${hash} ${message}`,
                    vscode.TreeItemCollapsibleState.None,
                    {
                        command: 'gitflow.log',
                        title: 'View Log',
                        arguments: []
                    },
                    '$(git-commit)'
                ));
            }
        } catch {
            // No commits or not a git repo
        }

        return items;
    }

    private getBranchItems(): GitFlowItem[] {
        const items: GitFlowItem[] = [];
        
        try {
            const result = child_process.execSync('git branch -a', { encoding: 'utf-8' });
            const lines = result.split('\n').filter(line => line.trim());

            for (const line of lines.slice(0, 10)) {
                const isCurrent = line.startsWith('*');
                const branch = line.substring(2).trim();
                
                items.push(new GitFlowItem(
                    branch,
                    vscode.TreeItemCollapsibleState.None,
                    {
                        command: 'gitflow.checkout',
                        title: 'Checkout',
                        arguments: [branch]
                    },
                    isCurrent ? '$(git-branch)' : '$(git-compare)'
                ));
            }
        } catch {
            // Not a git repo
        }

        return items;
    }

    private getRemoteItems(): GitFlowItem[] {
        const items: GitFlowItem[] = [];
        
        try {
            const result = child_process.execSync('git remote -v', { encoding: 'utf-8' });
            const lines = result.split('\n').filter(line => line.trim());
            const seen = new Set<string>();

            for (const line of lines) {
                const parts = line.split(/\s+/);
                if (parts.length >= 2) {
                    const name = parts[0];
                    const url = parts[1];
                    
                    if (!seen.has(name)) {
                        seen.add(name);
                        items.push(new GitFlowItem(
                            `${name}: ${url}`,
                            vscode.TreeItemCollapsibleState.None,
                            undefined,
                            '$(repo)'
                        ));
                    }
                }
            }
        } catch {
            // No remotes
        }

        return items;
    }

    private getStashItems(): GitFlowItem[] {
        const items: GitFlowItem[] = [];
        
        try {
            const result = child_process.execSync('git stash list', { encoding: 'utf-8' });
            const lines = result.split('\n').filter(line => line.trim());

            if (lines.length === 0) {
                items.push(new GitFlowItem(
                    'No stashes',
                    vscode.TreeItemCollapsibleState.None,
                    undefined,
                    '$(info)'
                ));
            } else {
                for (const line of lines.slice(0, 5)) {
                    items.push(new GitFlowItem(
                        line,
                        vscode.TreeItemCollapsibleState.None,
                        undefined,
                        '$(archive)'
                    ));
                }
            }
        } catch {
            // No stash
        }

        return items;
    }
}

export class GitFlowItem extends vscode.TreeItem {
    children: GitFlowItem[] | undefined;

    constructor(
        label: string,
        collapsibleState: vscode.TreeItemCollapsibleState,
        command?: vscode.Command,
        iconPath?: string
    ) {
        super(label, collapsibleState);
        this.command = command;
        if (iconPath) {
            this.iconPath = new vscode.ThemeIcon(iconPath.replace('$(', '').replace(')', ''));
        }
        this.contextValue = 'gitflowItem';
    }
}
