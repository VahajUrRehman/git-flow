import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { GitFlowTerminal } from './terminal';
import { GitFlowProvider } from './provider';
import { GitFlowStatusBar } from './statusbar';

let terminal: GitFlowTerminal | undefined;
let provider: GitFlowProvider | undefined;
let statusBar: GitFlowStatusBar | undefined;

export function activate(context: vscode.ExtensionContext) {
    console.log('GitFlow TUI extension is now active');

    // Initialize components
    terminal = new GitFlowTerminal();
    provider = new GitFlowProvider();
    statusBar = new GitFlowStatusBar();

    // Register tree data provider
    vscode.window.registerTreeDataProvider('gitflowView', provider);

    // Register commands
    const commands = [
        vscode.commands.registerCommand('gitflow.open', () => {
            terminal?.open();
        }),
        vscode.commands.registerCommand('gitflow.toggle', () => {
            terminal?.toggle();
        }),
        vscode.commands.registerCommand('gitflow.close', () => {
            terminal?.close();
        }),
        vscode.commands.registerCommand('gitflow.status', () => {
            terminal?.runCommand('status');
        }),
        vscode.commands.registerCommand('gitflow.log', () => {
            terminal?.runCommand('log --oneline -20');
        }),
        vscode.commands.registerCommand('gitflow.branch', () => {
            terminal?.runCommand('branch -a');
        }),
        vscode.commands.registerCommand('gitflow.commit', async () => {
            const message = await vscode.window.showInputBox({
                prompt: 'Enter commit message',
                placeHolder: 'Your commit message...'
            });
            if (message) {
                terminal?.runCommand(`commit -m "${message}"`);
            }
        }),
        vscode.commands.registerCommand('gitflow.push', () => {
            terminal?.runCommand('push');
        }),
        vscode.commands.registerCommand('gitflow.pull', () => {
            terminal?.runCommand('pull');
        }),
        vscode.commands.registerCommand('gitflow.stash', () => {
            terminal?.runCommand('stash');
        }),
        vscode.commands.registerCommand('gitflow.checkout', async () => {
            const branch = await vscode.window.showInputBox({
                prompt: 'Enter branch name',
                placeHolder: 'branch-name'
            });
            if (branch) {
                terminal?.runCommand(`checkout ${branch}`);
            }
        }),
        vscode.commands.registerCommand('gitflow.refresh', () => {
            provider?.refresh();
            statusBar?.update();
        }),
    ];

    // Add all commands to subscriptions
    commands.forEach(cmd => context.subscriptions.push(cmd));

    // Register status bar
    statusBar.register(context);

    // Watch for configuration changes
    vscode.workspace.onDidChangeConfiguration(e => {
        if (e.affectsConfiguration('gitflow')) {
            terminal?.dispose();
            terminal = new GitFlowTerminal();
        }
    });

    // Watch for git changes
    const gitExtension = vscode.extensions.getExtension('vscode.git');
    if (gitExtension) {
        gitExtension.activate().then(git => {
            const api = git.exports.getAPI(1);
            api.onDidChangeState(() => {
                provider?.refresh();
                statusBar?.update();
            });
        });
    }
}

export function deactivate() {
    terminal?.dispose();
    statusBar?.dispose();
}
