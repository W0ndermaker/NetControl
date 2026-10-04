# NetControl

Терминальная утилита для анализа, мониторинга и управления локальной сетью.

NetControl обнаруживает устройства в сети, определяет IP/MAC/hostname/производителя, сканирует порты, ведёт учёт адресов (IPAM) и выполняет ограниченные административные операции на удалённых Windows и Linux машинах по SSH.

**Статус:** 🟡 In Development · CLI first

Лицензия: [MIT](LICENSE)

---

## Подход

Проект начинается как **одна консольная утилита**, а не как веб-платформа. Всё, что делает веб-версия, сначала появляется в терминале: сканирование, мониторинг, IPAM, управление удалёнными машинами. Веб-интерфейс добавляется позже поверх тех же внутренних пакетов, CLI при этом остаётся полноценным клиентом.

Причины: быстрый отладочный цикл, работа по SSH без браузера, возможность использовать пайпы и скрипты, отсутствие инфраструктуры на старте.

Один бинарник, встроенное хранилище, ноль внешних сервисов.

## Возможности

- Обнаружение устройств в локальной сети (ARP → ICMP → TCP)
- Определение IP, MAC, hostname, производителя и ОС
- Сканирование портов
- Непрерывный мониторинг доступности и задержки в реальном времени
- IPAM: резервирование и освобождение IP, обнаружение конфликтов
- Удалённое управление Windows и Linux по SSH
- Временная блокировка Интернета с сохранением LAN и SSH
- Exam Mode для группы компьютеров
- История событий и журнал административных действий
- Вывод в виде таблицы или JSON для скриптов

## Установка

```bash
git clone <repository-url>
cd netcontrol
go build -o netcontrol ./cmd/netcontrol
```

Установка в `PATH`:

```bash
go install ./cmd/netcontrol
```

Требуется Go 1.22+. Для ARP и ICMP используются raw sockets, поэтому на Linux нужны `CAP_NET_RAW` (root или `setcap cap_net_raw+ep`).

## Использование

```text
netcontrol scan            сканировать сеть
netcontrol devices         список устройств
netcontrol device <ip>     подробности об устройстве
netcontrol ports <ip>      открытые порты
netcontrol net detect      определить параметры текущей сети
netcontrol watch           мониторинг в реальном времени
netcontrol ipam            работа с адресами
netcontrol ssh <host>      удалённое управление
netcontrol internet <host> блокировка и разблокировка Интернета
netcontrol exam            Exam Mode для группы компьютеров
netcontrol events          журнал событий
netcontrol config          конфигурация
netcontrol version         версия
```

### Сканирование

```bash
netcontrol scan                          # автоопределение сети
netcontrol scan --cidr 192.168.1.0/24   # конкретная сеть
netcontrol scan --iface enp3s0          # через интерфейс
netcontrol scan --ports 22,80,445       # свой список портов
netcontrol scan --deep                  # только TCP без discovery
```

Сканирование двухэтапное: сначала быстрый поиск активных хостов через ARP и ICMP, затем для найденных устройств проверяются порты, hostname и производитель. Порты по умолчанию — `22`, `53`, `80`, `443`, `445`, `3389`, `8080`, список конфигурируемый. Скан всех 65535 портов не входит в первую версию.

### Пример вывода

```text
$ netcontrol scan --cidr 192.168.1.0/24

Scanning 192.168.1.0/24 (254 hosts) with 50 workers

IP             MAC                 HOSTNAME   VENDOR       OS       LATENCY  STATUS
192.168.1.1    00:11:22:33:44:01   router     TP-Link      -        0.8 ms   online
192.168.1.21   aa:bb:cc:dd:ee:02   PC-01      Dell         Windows  1.4 ms   online
192.168.1.22   aa:bb:cc:dd:ee:03   debian     Supermicro   Linux    1.1 ms   online
192.168.1.50   aa:bb:cc:dd:ee:04   -          -            unknown  -        offline

4 devices, 3 online, 12 open ports, 1.9s
```

Прогресс и найденные устройства выводятся сразу, `Ctrl+C` корректно останавливает скан и освобождает ресурсы.

### Фильтры и вывод

```bash
netcontrol devices --status online
netcontrol devices --vendor Dell --since 24h
netcontrol devices --sort latency
netcontrol devices --json            # для jq и скриптов
netcontrol watch --interval 30s      # обновление по Ctrl+C
netcontrol watch --diff              # только изменения
```

Глобальные флаги: `--json`, `--no-color`, `--config <path>`, `--verbose`.

Коды возврата: `0` — успех, `1` — ошибка, `2` — частичный результат (часть устройств недоступна).

### IPAM

```bash
netcontrol ipam list --status free
netcontrol ipam reserve 192.168.1.50 --device PC-05 --comment "принтер"
netcontrol ipam release 192.168.1.50
netcontrol ipam check 192.168.1.50
```

Статусы IP: `USED`, `FREE`, `RESERVED`, `UNKNOWN`, `CONFLICT`. Перед резервированием адрес проверяется на занятость. Если один IP отвечает с двух разных MAC — фиксируется конфликт и IP получает статус `CONFLICT`.

> `FREE` означает, что NetControl не обнаружил устройство на этом IP. Это не гарантирует отсутствие DHCP lease.

### Удалённое управление

```bash
netcontrol ssh 192.168.1.21 test
netcontrol ssh 192.168.1.21 info
netcontrol ssh 192.168.1.21 firewall status

netcontrol internet 192.168.1.21 off --ttl 2h
netcontrol internet 192.168.1.21 on
```

Одинаковые команды для Windows и Linux — утилита сама определяет ОС и выбирает механизм.

| ОС | Стек | Действия |
| --- | --- | --- |
| Linux | SSH → nftables | test, info, firewall status, internet on/off |
| Windows | SSH → PowerShell → Defender Firewall | test, info, firewall status, internet on/off |

Блокировка Интернета сохраняет LAN, SSH и связь с самим NetControl. Сетевой интерфейс удалённой машины не выключается — это может разорвать SSH-соединение.

Каждая блокировка имеет safe rollback: по истечении TTL доступ восстанавливается автоматически, даже если NetControl уже не запущен.

### Exam Mode

```bash
netcontrol exam start --ip 192.168.1.21,192.168.1.22 --ttl 3h
netcontrol exam status
netcontrol exam stop
```

Блокирует Интернет на выбранной группе компьютеров, показывает прогресс и ошибки, затем восстанавливает доступ.

## Конфигурация

Файл `~/.config/netcontrol/config.yaml` (создаётся `netcontrol config init`), переменные окружения имеют приоритет.

| Переменная | Описание | По умолчанию |
| --- | --- | --- |
| `NETCONTROL_CIDR` | сеть для сканирования | автоопределение |
| `NETCONTROL_IFACE` | сетевой интерфейс | автоопределение |
| `NETCONTROL_PORTS` | порты для проверки | `22,53,80,443,445,3389,8080` |
| `NETCONTROL_WORKERS` | размер worker pool | `50` |
| `NETCONTROL_SCAN_TIMEOUT` | таймаут скана | `30s` |
| `NETCONTROL_WATCH_INTERVAL` | период мониторинга | `30s` |
| `NETCONTROL_SSH_TIMEOUT` | таймаут SSH | `10s` |
| `NETCONTROL_SSH_USER` | пользователь для SSH | текущий |
| `NETCONTROL_DB` | путь к хранилищу | `~/.local/share/netcontrol/netcontrol.db` |

Учётные данные SSH берутся из `~/.ssh/config` и ssh-agent либо задаются через `NETCONTROL_SSH_KEY`. Пароли и ключи не хранятся в репозитории.

## Архитектура

```text
cmd/netcontrol/main.go
        │
     cli/         команды, флаги, вывод
        │
    service/      логика: scan, ipam, watch, remote
        │
   scanner/  ssh/  firewall/  ipam/
        │
   store/         встроенное хранилище (SQLite)
```

Модульный монолит, один бинарник. Взаимодействие с сетью изолировано в `scanner` и `ssh`, поэтому те же пакеты переиспользуются веб-версией без изменений.

Хранилище встроенное (SQLite) — не требует запущенного сервера. PostgreSQL появится вместе с веб-версией.

## Структура проекта

```text
netcontrol/
├── cmd/netcontrol/main.go
├── internal/
│   ├── cli/           # команды, флаги, вывод
│   ├── config/
│   ├── service/
│   ├── scanner/       # arp, icmp, tcp, hostname, vendor
│   ├── store/         # репозитории и миграции
│   ├── ssh/           # client, connection, actions
│   ├── firewall/      # nftables, windows
│   ├── ipam/
│   └── output/        # table, json, color
├── migrations/
├── testdata/
├── go.mod
└── README.md
```

## Принципы

- Сначала рабочий CLI, потом веб — не наоборот.
- Никакого прочего shell-доступа: только перечисленные операции.
- Каждое изменение сети должно быть обратимым.
- Никакой телеметрии, внешних сервисов и лишних зависимостей.
- Один законченный проект лучше десяти недоделанных.

## Безопасность

- Произвольное выполнение shell-команд невозможно, доступен фиксированный набор действий.
- Административные операции требуют явного вызова и логируются в audit log: user, action, target, time, result.
- Блокировка Интернета всегда сохраняет доступ к управлению машиной.
- Изменения сети откатываются автоматически по TTL.
- SSH-учётные данные и секреты не попадают в исходный код и репозиторий.
- Резервирование IP не перезаписывает существующую конфигурацию сети.
- Фактические DHCP lease остаются источником истины.

## Roadmap

### CLI

| Phase | Содержание | Статус |
| --- | --- | --- |
| 0 | Repo, `go.mod`, структура, README, первый commit | ⬜ |
| 1 | Каркас CLI, команды, `version`, справка | ⬜ |
| 2 | Config: файл + env, `config init` | ⬜ |
| 3 | Network detection: auto + ручной CIDR | ⬜ |
| 4 | CIDR engine: network, broadcast, host range | ⬜ |
| 5 | ARP scanner | ⬜ |
| 6 | ICMP scanner и задержка | ⬜ |
| 7 | TCP scanner и список портов | ⬜ |
| 8 | Hostname и MAC OUI vendor | ⬜ |
| 9 | Worker pool, graceful shutdown | ⬜ |
| 10 | Context, timeout, отмена по Ctrl+C | ⬜ |
| 11 | Вывод: таблица, JSON, цвет, прогресс | ⬜ |
| 12 | Хранилище и реестр устройств | ⬜ |
| 13 | `devices`, `device`, `ports`, фильтры | ⬜ |
| 14 | `watch`: мониторинг и события online/offline | ⬜ |
| 15 | История устройства и журнал событий | ⬜ |
| 16 | IPAM: list, reserve, release, check | ⬜ |
| 17 | SSH client, аутентификация, test/info | ⬜ |
| 18 | Linux: nftables, internet on/off | ⬜ |
| 19 | Windows: PowerShell, Defender Firewall | ⬜ |
| 20 | Unified actions и safe rollback по TTL | ⬜ |
| 21 | Exam Mode | ⬜ |
| 22 | Тесты: unit + integration | ⬜ |
| 23 | Shell completions, man page, релиз v0.1.0 | ⬜ |

### Web (после CLI)

| Этап | Содержание |
| --- | --- |
| W1 | HTTP server, REST API поверх существующих сервисов |
| W2 | PostgreSQL вместо встроенного хранилища |
| W3 | WebSocket и real-time обновления |
| W4 | Frontend: dashboard, devices, IPAM, events |
| W5 | Auth, роли Admin / Viewer, audit log |

Легенда: ⬜ Not Started · 🟡 In Progress · 🟢 Completed · 🔴 Blocked

Этапы выполняются последовательно: текущий этап должен работать до начала следующего.
