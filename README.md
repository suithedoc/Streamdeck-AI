
# Streamdeck-AI

Streamdeck-AI is a fun side project written in Go that allows you to experiment with OpenAI functions using the Elgato Stream Deck. Please note that this is a hobby project without a focus on clean code and that everyone can use it at their own risk.

## Disclaimer: 
This is a hobby project and was not created with a focus on clean code. Use it at your own risk.

## Dependencies

To use this project, you need the following:

- Elgato Stream Deck
- OpenAI ApiKey
- Discord ApiKey (optional)

## Features
Streamdeck-AI has the following features:

- Text to speech integration (TTS)
- Use Stream Deck to activate bots
- Use command line to run commands or talk to Assistant
- Programmatically paste API response using the clipboard (e.g. into Minecraft console, Discord, etc.)
- Simple INI file to configure prompts
- Predefined bots for WhatsApp, Discord, Minecraft, Whisper, Chat Assistant, and Command-line Assistant

## Bots
Streamdeck-AI has several bots that you can activate using the Stream Deck or command line
1. Commander Bot
  - This bot allows you to execute Linux commands in the terminal.
  - Activated by Stream Deck button 0 (without history), button 5 (with history), and terminal input (multiline input, sending by pressing return twice).
  - Please note that this bot will try to execute the command directly, so use it at your own risk.
2. Assistant Bot
  - This bot is like Alexa, Siri, or ChatGPT but with shorter answers, and the response will be read out loud with TTS.
  - Activated by Stream Deck button 1 (without history), button 6 (with history), and terminal input prefix "a:" (without history) or "ah:" (with history).
3. Minecraft Bot
  - This bot allows you to execute Minecraft chat commands in the command prompt.
  - You can also have OpenAI create items or build stuff for you (which is not perfect but really fun!).
  - Activated by Stream Deck button 3.
  - Works by programmatically copying the response of OpenAI (GTP3.5/4) to the clipboard and pasting it to the Minecraft console.
4. Whisper Bot
  - This bot records your voice using Whisper, sends it to OpenAI-Whisper, copies it to the system clipboard, and pastes it wherever the cursor currently is.
  - Activated by Stream Deck button 2.
5. WhatsApp Bot
  - This bot allows OpenAI to answer WhatsApp messages.
  - Currently, one number and one name are configurable in the config.ini file.
6. Discord Bot
  - This bot allows you to have the power of OpenAI in a Discord channel.
7. Coding Bot
  - This bot allows OpenAI to create code for you and paste it directly into your IDE (to the cursor position).
  - Activated by Stream Deck button 4.

## Tips when having problems with the Stream Deck (tested on Manjaro Linux)
* create a udev rule:
  * `sudo nano /etc/udev/rules.d/99-streamdeck.rules`
* add following text into the .rules file
  * `SUBSYSTEMS=="usb", ATTRS{idVendor}=="vendor_id", ATTRS{idProduct}=="product_id", MODE="0660", TAG+="uaccess"`
* reboot

## Example config.ini file
```ini
apiKey=...
;model =gpt4
model = gpt3.5
whatsapp = disabled
discord = disabled
discordBotToken =
commanderSystemMsg= You are a Linux teacher that knows how to do everything on a Linux Arch installation with GNOME, with the command line only. You are eager to prove that you know a single line command for every request I give you.
commanderPromptMsg = Only answer in precise executable commands and prefer software that you know are installed. For emails, use Firefox. For web, use Firefox. Always execute everything in one line. Only show me commands that I could run as is, without edit. What is the command line for this request:
assistantSystemMsg = You are a personal assistant just like Alexa or Siri. Answer in short, precise sentences. For emails, use Firefox. For web, use Firefox. Always execute everything in one line.
assistantPromptMsg =
whatsappSystemMsg = You are Yoda from the Star Wars universe. Peter is your student. When you receive a message, it's from Peter. Talk to him as Yoda would. Use short sentences.
whatsappPromptMsg = Reply to Peter's message as Yoda and stay in character. Message:
whatsappNumber =
whatsappName =
minecraftSystemMsg = You are a Minecraft player who knows all the console commands when cheats are enabled. Try to create everything that is requested with console commands. Answer in executable Minecraft console commands only.
minecraftPromptMsg = Try to fulfill the request with Minecraft console commands only. Refer to myself as @s. Always put the command in code blocks, either one backtick for one line or three for multiple lines. Do not add comments to the code so I can straight copy and paste it to the console. Request:
discordSystemMsg = You are Yoda from the Star Wars universe. Peter is your student. When you receive a message, it's from Peter. Talk to him as Yoda would. Use short sentences.
discordPromptMsg = Reply to Peter's message as Yoda and stay in character. Message:
```

## Demo
Check out our demo videos for more information on how to use Streamdeck-AI.

## Contributing
We welcome all contributions, whether they are bug reports, feature requests, or pull requests. Please feel free to contribute!

## License
Streamdeck-AI is released under the [MIT License](https://opensource.org/licenses/MIT).
