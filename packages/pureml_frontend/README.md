[![PureMLUI](/assets/FrontendCoverImg.png)](https://pureml.com)

<br/>
<br/>

<div align="center">
  <a
    href="https://pypi.org/project/pureml/"
  >
    <img alt="Remix" src="https://img.shields.io/badge/remix-1.11.1-green?style=flat&logo=remix&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://python-poetry.org/"
  >
    <img src="https://img.shields.io/badge/react-18.2.0-red?style=flat&logo=react&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://opensource.org/licenses/Apache-2.0"
  >
    <img alt="License" src="https://img.shields.io/badge/tailwindcss-3.1.2-blue?style=flat&logo=tailwindcss&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://discord.gg/xNUHt9yguJ"
  >
    <img alt="Discord" src="https://img.shields.io/badge/Discord-Join%20Discord-blueviolet?style=flat&logo=discord&logoColor=white" />
  </a>
  &nbsp;
</div>

<br/>

## Quick start

PureML UI helps you to visualize all the details and versions of your models and datasets you want to work with. It takes no time to run PureML UI on your local system. Follow below steps to run:

Start the server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) to use the UI.

<br/>

## Directory Structure

```
pureml_frontend
├─ app/
│  ├─ analytics/             # Analytics
│  │  ├─ gtags.client.ts     # Google Tag
│  │  ├─ reportWebVitals.ts  # Vercel web vitals
│  │  └─ vitals.ts           # Vitals
|  |
|  ├─ components/            # Components
│  │  ├─ landingPage/        # Landing Page files
│  │  ├─ ui/                 # Minimal UI components used in all
│  │  └─ ...                 # other components used
│  │
│  ├─ routes/                # Route pages
│  │  ├─ api/                # All apis used in app
│  │  ├─ auth/               # Pages under Authentication
│  │  ├─ datasets/           # Pages showing all datasets
│  │  ├─ models/             # Pages showing all models
│  │  ├─ org/                # Pages under Org feature
│  │  ├─ settings/           # Pages under settings
│  │  └─ ...                 # …that has layouts of each page
│  │
│  ├─ entry.client.ts        # Entry.client file by remix
│  ├─ entry.server.ts        # Entry.server file by remix
│  ├─ lib.type.d.ts          # Type file
│  ├─ session.ts             # Session used for authentication
│  └─ root.tsx               # Root Index
│
├─ public/                   # Public
|  ├─ error/                 # Error images
|  └─ imgs/                  # Images used in application
|
├─ styles/                   # Styles
│  └─ app.css                # CSS file
|
├─ .gitignore                # List of files and folders not tracked by Git
├─ .eslintrc                 # Linting preferences for JavasScript
├─ remix.config.js           # Remix configuration file
├─ tailwind.config.js        # Tailwind configuration file
├─ package.json              # Project manifest
└─ README.md                 # This file
```

## Technology used

1. [Remix framework](https://remix.run/)
2. [Tailwind CSS](https://remix.run/)
3. [Tailwind Variants](https://www.tailwind-variants.org/docs/introduction)
4. [Radix components](https://www.radix-ui.com/docs/primitives/overview/introduction)
5. [Reactflow](https://reactflow.dev/)

## Reporting Bugs

To report any bugs you have faced while using PureML package, please

1. Report it in [Discord](https://discord.gg/xNUHt9yguJ) channel
2. Open an [issue](https://github.com/PureMLHQ/PureML/issues)

<br />

## Contributing and Developing

Lets work together to improve the features for everyone. Here's step one for you to go through our [Contributing Guide](./CONTRIBUTING.md). We are already waiting for amazing ideas and features which you all have got.

Work with mutual respect. Please take a look at our public [Roadmap here](https://pureml.notion.site/7de13568835a4cf18913307503a2cdd4?v=82199f96833a48e5907023c8a8d565c6).

<br />





# Windows

## Installation

To get started with Pureml, you'll need to install the PureML Python SDK. This SDK will allow you to work with the PureML platform, from creating your account to deploying your models.

To install PureML, simply run the following command:

```bash
pip3 install pureml
```

## New to PureML?
If you’re new to PureML, you’ll need to sign up for an account before you can start using the platform. CLI supports you to signup. You can do this using the pureml command line utility:

```bash
pureml auth signup
```
This command will prompt you to enter your email, user handle, name, and password to create your PureML account. Once you register your details successfully, you will receive a verification mail on registered email Id to verify and proceed for login.

## Sign In to start what you left
If you already have a PureML account, you can log in using 3 different ways:

### 1.CLI
The `pureml` command line utility:

```bash
pureml auth login
```

### 2.SDK
You can sign in to PureML using our Python SDK:
```bash
import pureml

pureml.login(org_id='<org_id>', access_token='<access_token>')
```

### 3. PureML UI
For cloud users, go to `https://pureml.com/auth/signin` to sign in through your sign in credentials. For self hosted users, go to `http://localhost:3000/auth/signin`

Once you’re logged in, you’ll be able to view your datasets, models, and other assets using the PureML platform.

# Self-host PureML

If you`re a large company looking for a proof-of-concept, or an engineer looking to use the open-source version in a funky non-production way, then you’re in the right place! Our Docker compose deployment let’s you spin up a fresh Pureml instance in minutes.

Want more reliability? The easiest way to get started with PureML is to use [PureML Cloud](https://pureml.com/)

## Requirements
A subsystem with docker-engine installed. [Installation guidelines](https://docs.docker.com/engine/install/)

## Setting up the compose
1. In a new directory where you want to setup the containers, create a new file `docker-compose.yml`
2. Add the following content from our [official docker-compose example file](https://raw.githubusercontent.com/PuremlHQ/PureML/main/packages/pureml_docker/docker-compose.yml)

```bash
version: "3"

services:
  backend:
    image: puremlhq/pureml_backend:local-base
    environment:
      - PURE_SITE_BASE_URL=http://localhost:3000
    ports:
      - 8080:8080
    volumes:
      - pureml-data:/pureml_backend/data

  frontend:
    image: puremlhq/pureml_frontend:dev
    environment:
      - NEXT_PUBLIC_BACKEND_URL=http://backend:8080/api/
    ports:
      - 3000:3000
    links:
      - backend

volumes:
  pureml-data:

```
3. Run the following command to start your containers `docker compose up`
Make sure your docker engine is running for the docker command to work.

4. Additionally, to run the containers in background you can use the command `docker compose up -d` or `docker compose up --detach`

If all goes well, You should have a PureML local instance setup and running PureML UI at [localhost:3000](https://localhost:3000/) with your backend hosted at [localhost:8080/api](http://localhost:8080/api). You can even checkout the auto generated Open API swagger documentation at [/api/swagger/index.html](http://localhost:8080/api/swagger/index.html)

## Usage
You can use your local self-hosted backend in the sdk by using the `set_backend` function Example:

Assuming your backend url is [localhost:8080/api](http://localhost:8080/api)

```bash
import pureml

# Call this function near the initialization of your scripts, preferably just after importing
pureml.set_backend("http://localhost:8080/api/") # trainling slash at the end is important

...
```
## Configuration
There are various ways to configure and personalize your PureML instance to better suit your needs. In this section you will find all the information you need about settings and options you can configure to get what you need out of PureML.


## Environment variables
### Frontend environment
| Variable | Variable | Variable |
| :------- | :------- | :------- |
|`NEXT_PUBLIC_BACKEND_URL`          |backend url env to connect frontend          |    [http://backend:8080/api/](http://backend:8080/api/)      |


## Backend environment
There are multiple environment variables that can be configured externally. But the application also depends upon the configuration of the docker image used to build the compose.

For Base Configuration : Rows with a missing ‘Default Value’ usually default to an empty string. This is different from `None`.

|  Variable     |   Description                         |  	Default Value           | 	
| :---- | :------------------------- | :---------- |       
| PORT     |	port where backend will be deployed                            | 	8080             |
|DATABASE       |type of database which will be used. Currently supported sqlite3 & postgres                            |sqlite3             |
|DATABASE_URL       |url for the database (❗Required if database type is postgres)                            |             |
|CGO_ENABLED       |	whether to use cgo (requires gcc) or pure go                            |1             |
|PURE_SITE_BASE_URL       |frontend url for email redirections                            |             |





## Community

To get quick updates of feature releases of PureML, follow us on:

[<img alt="Twitter" height="20" src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" />](https://twitter.com/getPureML) [<img alt="LinkedIn" height="20" src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" />](https://www.linkedin.com/company/PuremlHQ/) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white" />](https://github.com/PureMLHQ/PureML) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" />](https://discord.gg/DBvedzGu)

<br/>

## 📄 License

See the [Apache-2.0](./License) file for licensing information.
