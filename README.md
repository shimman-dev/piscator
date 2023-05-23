# piscator - fishing made easy

![intro text and piscator demo](./docs/demo.gif)

Ahoy, fellow sailor! Set your sights on the `piscator` CLT, a trusty companion
for your coding odysseys. With `piscator` by your side, you'll navigate the vast
seas of GitHub with ease, harnessing the power of the command-line interface to
explore, capture, and conquer the realm of repositories.

As a seasoned sailor, you'll command the `piscator` vessel to cast its net
into the GitHub sea, deftly capturing the URLs of repositories from users,
organizations, or even yourself. With each cast, `piscator` will bring you a
bountiful haul of valuable code treasures, presenting them in a convenient
format for further exploration. But `piscator`'s abilities don't end there! With
the reel command, you'll transform those captured repositories into tangible
assets by creating dedicated directories and git repositories for each catch,
ready to be harnessed for your coding voyages.

With its intuitive commands, `piscator` empowers you to hoist the sails of
productivity and embark on coding adventures like never before. Whether you're
seeking to explore, collect, or reel in repositories, `piscator` will be your
loyal first mate, guiding you through the GitHub waters with precision and
efficiency. Embrace the spirit of the sailor, unleash the power of
`piscator`, and embark on a coding voyage that will leave a lasting mark on the
open seas of code.

---

Table of contents:

- [Installation](#installation)
- [Usage](#usage)
- [Documentation](#documentation)
- [Todos](#todos)

---

## [Installation](#installation)

Although the current [releases](https://github.com/shimman-dev/piscator/releases)
page offers binaries for Linux, MacOS, and Windows, rest assured that I have
heard your concerns and I am committed to providing ye with the convenience and
ease of installation ye deserve. I solemnly swear to make future
[releases](https://github.com/shimman-dev/piscator/releases) available through
all major package managers, for I believe in the principle that good things
come to sailors who trust their captain.

So, my dear crew, hold tight to the mast and keep your hearts filled with hope.
In the near future, ye shall witness the expansion of piscator's reach, as it
sets sail on the vast seas of package managers, bringing ye the treasures of
effortless installation and seamless updates. Trust in me, your captain, for I
am dedicated to improving the voyage of piscator and ensuring that every sailor
can embark on their coding adventures with the wind at their backs.

Together, we shall overcome this minor setback and forge ahead toward a
brighter future for piscator. Ye are the lifeblood of this crew, and I am
grateful for your unwavering support. May our sails be filled with the winds of
progress, as we navigate the ever-changing tides of technology. Onward, my
brave crew, for the best is yet to come!

### [Usage](#usage)

Avast, matey! Let me tell ye about the heart of the `piscator ` vessel: the
mighty commands `cast` and `reel`.

With a swift cast, the command will fetch ye a bounty of JSON, ripe for the
takin'. This JSON treasure can be piped into other commands like `fzf` or `jq`
to uncover its secrets. And if ye be desirin' a written record of yer findings,
simply wave the `repos.json` flag (`-f` or `--makeFile`), and a file named
`repos.json` shall be fashioned right where ye stand (directory where `piscator`
was ran).

But behold, the power of the mighty reel! This command harnesses the fruits of
`cast` to weave a directory in the likeness of the chosen user or organization.
Then, with the skill of an expert angler, it reels in each specified repository,
be it from the depths of the sea or the local waters. Should a repository
already rest in the directory, fear not! `reel` will deftly pull the latest
catch, ensuring ye always have the freshest head.

So, me hearty, set sail with `piscator` and let the commands `cast` and `reel`
guide ye on yer voyage to GitHub riches!

## [Documentation](#documentation)

### [piscator](#piscator)

Running `piscator help` will show the commands of `piscator`. You can also run
`piscator cast -h` and `piscator reel -h` to see the flags each respective
command takes.

![running piscator help](./docs/piscator-help.gif)

### [cast](#cast)

Running `piscator cast username` will output a JSON of public repositories
for a user:

![running piscator cast shieldbattery](./docs/cast-user.gif)

<details>
<summary>example output of `piscator cast shieldbattery`</summary>

```text
[
{
	"name": "broodmap",
	"html_url": "https://github.com/ShieldBattery/broodmap",
	"language": "Rust",
	"fork": false,
	"private": false,
	"size": 4695
},
{
	"name": "rally-point",
	"html_url": "https://github.com/ShieldBattery/rally-point",
	"language": "JavaScript",
	"fork": false,
	"private": false,
	"size": 966
},
{
	"name": "ShieldBattery",
	"html_url": "https://github.com/ShieldBattery/ShieldBattery",
	"language": "TypeScript",
	"fork": false,
	"private": false,
	"size": 244880
},
{
	"name": "stimpack",
	"html_url": "https://github.com/ShieldBattery/stimpack",
	"language": "Rust",
	"fork": false,
	"private": false,
	"size": 57
},
{
	"name": "trrr",
	"html_url": "https://github.com/ShieldBattery/trrr",
	"language": "Rust",
	"fork": false,
	"private": false,
	"size": 18
}
]
```

</details>

---

Running `piscator cast your_username -s` will output a JSON of your repositories
(public and private; requires a personal access token).

![running piscator cast azemetre -s](./docs/cast-self.gif)

The `-s` flag refers to your self, i.e. your GitHub username. The `-t` flag
takes a [PAT (how to create a PAT)](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token), alternatively you can export your PAT in your terminal environment with
the following command (this is demonstrated in the gif above):

`export GITHUB_TOKEN=pat_token_here`

Please note, that if you update your PAT you may need to remove your current
exported token with the following commands:

```
// remove token
unset GITHUB_TOKEN

// verify token disappeared
printenv GITHUB_TOKEN
```

<details>
<summary>example output of `piscator cast azemetre -s`</summary>

```text
[
{
	"name": "auteur-palettes",
	"html_url": "https://github.com/azemetre/auteur-palettes",
	"language": "JavaScript",
	"fork": false,
	"private": true,
	"size": 403
},
{
	"name": "azemetredotcom",
	"html_url": "https://github.com/azemetre/azemetredotcom",
	"language": "JavaScript",
	"fork": false,
	"private": true,
	"size": 24256
},
{
	"name": "hipster.nvim",
	"html_url": "https://github.com/azemetre/hipster.nvim",
	"language": "Lua",
	"fork": false,
	"private": false,
	"size": 1446
},
{
	"name": "musical-adventure",
	"html_url": "https://github.com/azemetre/musical-adventure",
	"language": "TypeScript",
	"fork": false,
	"private": true,
	"size": 993
},
{
	"name": "piscator",
	"html_url": "https://github.com/shimman-dev/piscator",
	"language": "Go",
	"fork": false,
	"private": false,
	"size": 1006
}
]
```

</details>

---

Running `piscator cast org_name -o` will output a JSON of public and repositories for an organization:

![running piscator cast shimman-dev -o](./docs/cast-org.gif)

**Please note:** as with the `-s` flag (`--self`), the `--org` requires a PAT
passed with the `--token` flag or fed into the env variable `GITHUB_TOKEN`.

<details>
<summary>example output of `piscator cast shimman-dev -o`</summary>

```text
[
{
	"name": "eslint-config",
	"html_url": "https://github.com/shimman-dev/eslint-config",
	"language": "JavaScript",
	"fork": false,
	"private": false,
	"size": 227
},
{
	"name": "piscator",
	"html_url": "https://github.com/shimman-dev/piscator",
	"language": "Go",
	"fork": false,
	"private": false,
	"size": 1006
},
{
	"name": "knockerupper",
	"html_url": "https://github.com/shimman-dev/knockerupper",
	"language": "",
	"fork": false,
	"private": true,
	"size": 14
}
]
```

</details>

---

Running `piscator cast username -x` will output a JSON of public and forked repositories:

![running piscator cast azemetre -x](./docs/cast-fork.gif)

<details>
<summary>example output of `piscator cast azemetre -x`</summary>

```text
[
{
	"name": "Adv360-Pro-ZMK",
	"html_url": "https://github.com/azemetre/Adv360-Pro-ZMK",
	"language": "",
	"fork": true,
	"private": false,
	"size": 145
},
	{
	"name": "auteur-palettes",
	"html_url": "https://github.com/azemetre/auteur-palettes",
	"language": "JavaScript",
	"fork": false,
	"private": true,
	"size": 403
}
]
```

</details>

---

Running `piscator cast username -f` will output a JSON of public repositories
and create a `repos.json` file:

![running piscator cast azemetre -x](./docs/cast-file.gif)

### [reel](#reels)

**Please note:** `piscator reel` can take the same flags as `piscator cast`, so
if you would like to reel in repos from an organization or yourself it will
require the same flags and arguments.

Running `piscator reel org_name` create a directory of the user/org and clone
their repositories:

![running piscator reel azemetre](./docs/reel-user.gif)

## [Todos](#todos)

- [x] flesh out readme
- [x] create vhs tapes
- [x] automate release binaries
- [ ] release on major package managers
- [ ] homebrew
- [ ] nix
- [ ] fedora
- [ ] macports
- [ ] arch linux (btw)
- [ ] debian
- [ ] scoop
- [x] add tests
- [ ] generate man pages
- [ ] make showcase site
- [ ] add ability to filter by language
