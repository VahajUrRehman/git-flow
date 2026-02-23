-- GitFlow TUI Neovim Plugin
-- A complete Git management interface for Neovim

local M = {}

-- Configuration
M.config = {
	-- Binary path (will be auto-detected if not set)
	bin_path = nil,
	-- Window settings
	window = {
		width = 0.9,
		height = 0.9,
		border = "rounded",
	},
	-- Keymaps
	keymaps = {
		open = "<leader>gg",
		close = "q",
	},
	-- Theme (uses Neovim colors if not set)
	theme = {
		primary = "#00D9A5",
		secondary = "#00B4A6",
		tertiary = "#0091EA",
		accent = "#00E5FF",
		highlight = "#FF6D00",
	},
}

-- State
local state = {
	buf = nil,		-- Terminal buffer
	win = nil,		-- Window handle
	job = nil,		-- Job ID
	is_open = false,
}

-- Find the gitflow-tui binary
local function find_binary()
	if M.config.bin_path then
		return M.config.bin_path
	end

	-- Try common locations
	local paths = {
		vim.fn.stdpath("data") .. "/gitflow-tui/gitflow-tui",
		vim.fn.exepath("gitflow-tui"),
		"./gitflow-tui",
		"gitflow-tui",
	}

	for _, path in ipairs(paths) do
		if vim.fn.executable(path) == 1 then
			return path
		end
	end

	return nil
end

-- Get project root (git repository root)
local function get_project_root()
	local handle = io.popen("git rev-parse --show-toplevel 2>/dev/null")
	if handle then
		local result = handle:read("*a")
		handle:close()
		return vim.trim(result)
	end
	return vim.fn.getcwd()
end

-- Calculate window dimensions
local function get_window_config()
	local width = math.floor(vim.o.columns * M.config.window.width)
	local height = math.floor(vim.o.lines * M.config.window.height)
	local col = math.floor((vim.o.columns - width) / 2)
	local row = math.floor((vim.o.lines - height) / 2)

	return {
		relative = "editor",
		width = width,
		height = height,
		col = col,
		row = row,
		style = "minimal",
		border = M.config.window.border,
	}
end

-- Open GitFlow TUI
function M.open()
	if state.is_open then
		M.focus()
		return
	end

	local bin = find_binary()
	if not bin then
		vim.notify("GitFlow TUI binary not found. Please install it.", vim.log.levels.ERROR)
		return
	end

	-- Create buffer
	state.buf = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_option(state.buf, "bufhidden", "hide")
	vim.api.nvim_buf_set_option(state.buf, "filetype", "gitflow")

	-- Create window
	local win_config = get_window_config()
	state.win = vim.api.nvim_open_win(state.buf, true, win_config)
	vim.api.nvim_win_set_option(state.win, "winhl", "Normal:Normal")

	-- Start terminal
	local cmd = bin .. " --cwd=" .. get_project_root()
	state.job = vim.fn.termopen(cmd, {
		on_exit = function(_, code, _)
			state.is_open = false
			if code == 0 then
				M.close()
			end
		end,
	})

	-- Set terminal options
	vim.cmd("startinsert")
	vim.api.nvim_buf_set_keymap(state.buf, "t", M.config.keymaps.close, 
		"<C-\\><C-n>:lua require('gitflow').close()<CR>", 
		{ noremap = true, silent = true })

	-- Set up autocommands for window resize
	vim.api.nvim_create_autocmd("VimResized", {
		callback = function()
			if state.is_open and state.win then
				vim.api.nvim_win_set_config(state.win, get_window_config())
			end
		end,
	})

	state.is_open = true
end

-- Close GitFlow TUI
function M.close()
	if state.win and vim.api.nvim_win_is_valid(state.win) then
		vim.api.nvim_win_close(state.win, true)
	end
	if state.buf and vim.api.nvim_buf_is_valid(state.buf) then
		vim.api.nvim_buf_delete(state.buf, { force = true })
	end
	state.win = nil
	state.buf = nil
	state.job = nil
	state.is_open = false
end

-- Focus the GitFlow window
function M.focus()
	if state.win and vim.api.nvim_win_is_valid(state.win) then
		vim.api.nvim_set_current_win(state.win)
		vim.cmd("startinsert")
	end
end

-- Toggle GitFlow TUI
function M.toggle()
	if state.is_open then
		M.close()
	else
		M.open()
	end
end

-- Setup function
function M.setup(opts)
	M.config = vim.tbl_deep_extend("force", M.config, opts or {})

	-- Create commands
	vim.api.nvim_create_user_command("GitFlow", M.open, {})
	vim.api.nvim_create_user_command("GitFlowToggle", M.toggle, {})
	vim.api.nvim_create_user_command("GitFlowClose", M.close, {})

	-- Set up keymaps
	if M.config.keymaps.open then
		vim.keymap.set("n", M.config.keymaps.open, M.toggle, 
			{ noremap = true, silent = true, desc = "Toggle GitFlow TUI" })
	end

	-- Set up highlights
	M.setup_highlights()
end

-- Setup highlight groups
function M.setup_highlights()
	local colors = M.config.theme

	vim.api.nvim_set_hl(0, "GitFlowPrimary", { fg = colors.primary, bold = true })
	vim.api.nvim_set_hl(0, "GitFlowSecondary", { fg = colors.secondary })
	vim.api.nvim_set_hl(0, "GitFlowTertiary", { fg = colors.tertiary })
	vim.api.nvim_set_hl(0, "GitFlowAccent", { fg = colors.accent })
	vim.api.nvim_set_hl(0, "GitFlowHighlight", { fg = colors.highlight })
end

-- Git command wrappers
function M.git_command(cmd)
	local bin = find_binary()
	if not bin then
		vim.notify("GitFlow TUI binary not found", vim.log.levels.ERROR)
		return
	end

	local handle = io.popen(bin .. " " .. cmd .. " 2>&1")
	if handle then
		local result = handle:read("*a")
		handle:close()
		return vim.trim(result)
	end
end

-- Quick actions
function M.status()
	M.open()
end

function M.log()
	M.git_command("log --oneline -20")
end

function M.branch()
	M.git_command("branch -a")
end

function M.commit(message)
	if message then
		M.git_command('commit -m "' .. message .. '"')
	else
		M.open()
	end
end

function M.push(remote, branch)
	remote = remote or "origin"
	branch = branch or ""
	M.git_command("push " .. remote .. " " .. branch)
end

function M.pull(remote, branch)
	remote = remote or "origin"
	branch = branch or ""
	M.git_command("pull " .. remote .. " " .. branch)
end

function M.checkout(branch)
	M.git_command("checkout " .. branch)
end

function M.stash()
	M.git_command("stash")
end

function M.stash_pop()
	M.git_command("stash pop")
end

-- Integration with other plugins
function M.telescope_integration()
	local has_telescope, telescope = pcall(require, "telescope")
	if not has_telescope then
		vim.notify("Telescope not found", vim.log.levels.WARN)
		return
	end

	local pickers = require("telescope.pickers")
	local finders = require("telescope.finders")
	local conf = require("telescope.config").values
	local actions = require("telescope.actions")
	local action_state = require("telescope.actions.state")

	local function gitflow_picker(opts)
		opts = opts or {}
		pickers.new(opts, {
			prompt_title = "GitFlow Actions",
			finder = finders.new_table({
				results = {
					{ "Open GitFlow TUI", M.open },
					{ "Git Status", function() M.git_command("status") end },
					{ "Git Log", function() M.git_command("log --oneline -20") end },
					{ "Git Branches", function() M.git_command("branch -a") end },
					{ "Git Stash List", function() M.git_command("stash list") end },
				},
				entry_maker = function(entry)
					return {
						value = entry,
						display = entry[1],
						ordinal = entry[1],
						action = entry[2],
					}
				end,
			}),
			sorter = conf.generic_sorter(opts),
			attach_mappings = function(prompt_bufnr, map)
				actions.select_default:replace(function()
					actions.close(prompt_bufnr)
					local selection = action_state.get_selected_entry()
					if selection then
						selection.action()
					end
				end)
				return true
			end,
		}):find()
	end

	-- Register with Telescope
	telescope.register_extension({
		exports = {
			gitflow = gitflow_picker,
		},
	})
end

-- Statusline component
function M.statusline()
	local branch = vim.trim(vim.fn.system("git branch --show-current 2>/dev/null") or "")
	if branch == "" then
		return ""
	end

	local status = vim.trim(vim.fn.system("git status --porcelain 2>/dev/null") or "")
	local indicator = ""
	if status ~= "" then
		indicator = " *"
	end

	return "  " .. branch .. indicator .. " "
end

return M
