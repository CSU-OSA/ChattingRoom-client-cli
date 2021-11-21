# ChattingRoom-client-cli

CSU-OSA Chatting Room (command line client)

CSU-OSA 聊天室（命令行客户端）

## Usage

### Command

- `/server http://localhost:8003/chat`
- `/user`
  - `login nick ticket`
  - `logout nick`
  - `siwtch nick`
- `/channel`
  - `create name ticket`
  - `join name ticket`

### .chattingroomrc

*配置文件的临时替代*

可在程序根目录下的`.chattingroomrc`文件中预先写入指令，客户端将会在启动时自动读取并运行。

## TODO

- [x] multi-user
- [ ] multi-channel
- [ ] better command
- [ ] configure file
- [ ] terminal UI
