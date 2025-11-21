# Установка Task (go-task) — кратко и понятно

Этот файл объясняет три способа установки Task CLI на macOS, Windows и Ubuntu:
1. Основной (рекомендуемый) — через curl (установочный скрипт / репозитории).
2. Второй — через npm (@go-task/cli).
3. Третий (альтернативный / удобный) — через Homebrew (macOS) / пакетный менеджер (apt/snap/choco/winget).

Проверка после установки:
```bash
task --version
```

--------------------------------------------------------------------------------
macOS
--------------------------------------------------------------------------------

1) Через curl (рекомендуемый — переносимый бинарник)
- Быстро, подходит для CI и локальной установки.
- Установить в `~/.local/bin`:
```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -b ~/.local/bin
```
- Добавьте в PATH (если ещё не добавлен):
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc   # или ~/.zshrc
source ~/.bashrc   # или source ~/.zshrc
```
- Проверка:
```bash
task --version
```

2) Через npm (второй вариант, кроссплатформенно)
- Удобно, если у вас уже установлен Node.js/npm.
```bash
npm install -g @go-task/cli
task --version
```
- На macOS глобальные npm-пакеты могут устанавливаться в каталоги, которые не в PATH — проверьте `npm root -g` и при необходимости добавьте в PATH.

3) Через Homebrew (альтернативный/посредственный)
- Если вы пользуетесь Homebrew:
```bash
brew install go-task/tap/go-task
# или
brew install go-task
task --version
```

--------------------------------------------------------------------------------
Windows
--------------------------------------------------------------------------------

Самый простой способ (рекомендуемый):
1) Через Winget (самый простой для Windows)
```powershell
winget install Task.Task
# затем в новой консоли
task --version
```

2) Через npm (альтернативно)
- Откройте PowerShell или CMD с правами пользователя:
```powershell
npm install -g @go-task/cli
task --version
```

3) Через Chocolatey или Scoop (ещё один вариант)
- Chocolatey:
```powershell
choco install go-task
task --version
```
- Scoop:
```powershell
scoop install task
task --version
```

Примечание: если используете WSL / Git Bash — можно также применить curl-скрипт из раздела Linux/macOS.

--------------------------------------------------------------------------------
Ubuntu (и Debian-подобные)
--------------------------------------------------------------------------------

1) Через curl + apt (рекомендуемый, простой для Ubuntu)
- Настроить репозиторий и установить пакет:
```bash
curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | sudo -E bash
sudo apt update
sudo apt install task
task --version
```

2) Через npm (альтернатива)
```bash
sudo npm install -g @go-task/cli   # или без sudo при корректной настройке npm
task --version
```

3) Через Snap (альтернативный способ)
```bash
sudo snap install task --classic
task --version
```

--------------------------------------------------------------------------------
Дополнительные замечания
--------------------------------------------------------------------------------

- Установка в пользовательскую директорию:
  - Скрипт установки поддерживает опции `-b` для указания директории установки:
    ```bash
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -b ~/.local/bin
    ```
  - Если поставили в пользовательскую директорию, не забудьте добавить её в PATH.

- Версии:
  - Чтобы поставить конкретную версию через скрипт:
    ```bash
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -b ~/.local/bin v3.42.1
    ```

- Автодополнения для оболочки:
  - После установки можно подключить автодополнение:
    ```bash
    # bash
    echo 'eval "$(task --completion bash)"' >> ~/.bashrc

    # zsh
    echo 'eval "$(task --completion zsh)"' >> ~/.zshrc
    ```

- Если установка через npm приводит к сообщению о правах доступа, либо используйте `sudo` (не всегда желательно) либо настройте глобальный префикс npm в домашней директории.

--------------------------------------------------------------------------------
Краткое резюме (что выбрать)
- macOS: curl (рекомендуется) → npm (если уже есть npm) → brew (если используете Homebrew).
- Windows: winget (самый простой) → npm → choco/scoop.
- Ubuntu: curl→apt script (рекомендуется) → npm → snap.
