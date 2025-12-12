# Debugging Triggers - Step by Step

I've added comprehensive debug logging to help identify why triggers aren't working.

## Steps to Debug

1. **Rebuild the application:**
   ```powershell
   .\scripts\setup-tdm-gcc.ps1
   go build main.go
   ```

2. **Run the application:**
   ```powershell
   .\main.exe
   ```
   (Or use `go run main.go`)

3. **Watch the terminal output** - you should see:
   - `keyboard hook started successfully`
   - `[DEBUG] KeyboardHook: starting keyboard hook...`
   - `[DEBUG] KeyboardHook: keyboard hook started, waiting for events...`

4. **Type some keys in Notepad** - you should see debug messages like:
   - `[DEBUG] KeyboardHook: received key event #1: ";" (rawcode=...`
   - `[DEBUG] OnKeyPress: received key: ";", buffer before: ""`
   - etc.

5. **Type a trigger** (e.g., `;email` + Space):
   - You should see the buffer being built up
   - When you press Space, you should see:
     - `[DEBUG] OnKeyPress: SPACE, checking expansion, buffer: ";email"`
     - `[DEBUG] CheckAndExpand: buffer content: ";email"`
     - `[DEBUG] CheckAndExpand: checking X expansions`
     - Either a match or "no matching expansion found"

6. **Check what you see:**
   - **If you see NO keyboard events**: The keyboard hook isn't receiving events (permissions issue?)
   - **If you see events but buffer is empty**: Buffer isn't being populated
   - **If you see buffer but no expansion check**: Settings might be disabled
   - **If you see expansion check but no match**: Trigger matching logic issue

## Common Issues Based on Debug Output

### No keyboard events at all
- **Problem**: Keyboard hook isn't receiving events
- **Solution**: Run as Administrator, check Windows permissions

### Keyboard events but callback is nil
- **Problem**: Callback wasn't set properly
- **Solution**: Check if expander.Start() was called

### Buffer not being populated
- **Problem**: Keys aren't being recognized as printable
- **Solution**: Check translateEvent function

### Expansion check not happening
- **Problem**: Settings disabled or trigger keys disabled
- **Solution**: Check config/expansions.json settings

### Expansion check happens but no match
- **Problem**: Trigger matching logic
- **Solution**: Check if trigger is in buffer exactly as defined

## After Debugging

Once we identify the issue, we can remove or reduce the debug logging for production use.

