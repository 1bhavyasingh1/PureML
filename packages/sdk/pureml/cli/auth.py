import socket
from time import sleep, time
import requests
import typer
from rich import print
from pureml.cli.helpers import save_auth

from pureml.schema.backend import get_backend_base_url, get_frontend_base_url
from .orgs import select
from rich.progress import Progress, SpinnerColumn, TextColumn
from pureml.components import delete_token, get_org_id, get_token
from urllib.parse import urljoin
import json
import platform
import ipapi

app = typer.Typer()

def get_location():
    try:
        response = ipapi.location(output="json")
    except:
        try:
            # print("Getting device details...")
            response = requests.get(f'https://api64.ipify.org?format=json').json()
            response = requests.get(f"https://ipapi.co/{response['ip']}/json/", headers={
                    "User-Agent": "pureml-cli"
                }).json()
        except:
            response = {
                "ip": socket.gethostbyname(socket.gethostname()),
                "city": "Unknown",
                "region": "Unknown",
                "country": "Unknown",
            }
    location_data = {
        "ip": response["ip"] or socket.gethostbyname(socket.gethostname()),
        "city": response["city"] or "Unknown",
        "region": response["region"] or "Unknown",
        "country": response["country_name"] or "Unknown",
    }
    return location_data

@app.command()
def details():
    token = get_token()
    org_id = get_org_id()

    print("Org Id: ", org_id)
    print("Access Token: ", token)


@app.callback()
def callback():
    """
    Authentication user command

    Use with status, signup, login or logout option

    status - Checks current authenticated user
    signup - Creates new user
    login - Logs in user
    logout - Logs out user
    """


@app.command()
def signup(backend_url: str = typer.Option("", "--backend-url", "-b", help="Backend URL for self-hosted or custom pureml backend instance")):
    try:
        print("\nCreate a new account/\n")
        email: str = typer.prompt("Enter new email")
        handle: str = typer.prompt("Enter new user handle")
        name: str = typer.prompt("Enter new user name")
        password: str = typer.prompt(
            "Enter new password", confirmation_prompt=True, hide_input=True
        )
        # organization_id: str = typer.prompt("Enter Organization id")
        # data = {"email": email, "password": password, "org": organization_id}
        data = {"email": email, "password": password, "handle": handle, "name": name}
        backend_base_url = get_backend_base_url(backend_url)

        url_path_1 = "user/signup"
        url = urljoin(backend_base_url, url_path_1)

        response = requests.post(url, json=data)

        if not response.ok:
            print(f"[red]Could not create account! Please try again later")
            return
        print(
            f"[green]Successfully created your account! You can now login using ```pure auth login```"
        )
    except Exception as e:
        print(f"[red]Could not create account! Please try again later")
        print(e)

@app.command()
def login(
    backend_url: str = typer.Option("", "--backend-url", "-b", help="Backend URL for self-hosted or custom pureml backend instance"),
    frontend_url: str = typer.Option("", "--frontend-url", "-f", help="Frontend URL for self-hosted or custom pureml frontend instance"),
    browserless: bool = typer.Option(False, "--browserless", "-s", help="Browserless login for ssh or pipelines"),
    interactive: bool = typer.Option(False, "--interactive", "-i", help="Login with email and password interactively"),
    ):
    try:
        backend_base_url = get_backend_base_url(backend_url)
        frontend_base_url = get_frontend_base_url(frontend_url)
        # Interactive login with email and password
        if interactive:
            print(f"\n[Enter your credentials to login[/\n")
            email: str = typer.prompt("Enter your email")
            password: str = typer.prompt("Enter your password", hide_input=True)
            data = {"email": email, "password": password}

            url_path = "user/login"
            url = urljoin(backend_base_url, url_path)

            response = requests.post(url, json=data)

            if response.ok:
                token = response.text
                token = json.loads(token)["data"][0]

                access_token = token["accessToken"]
                email = token["email"]

                save_auth(access_token=access_token, email=email)

                print(f"[green]Successfully logged in as {email}!")

                # Select organization
                org_id = select()
                if org_id is not None:
                    save_auth(org_id=org_id, access_token=access_token, email=email)
                    print(f"[green]Successfully linked to organization {org_id}!")
                else:
                    print(f"[red]Organization details not found! Please contact your admin to add you to an organization!")

            else:
                print(f"[red]Unable to login to your account!")
        
        # Browser based login
        else:
            # Get device details
            device = platform.platform()
            device_data = get_location()
            device_id = device_data["ip"]
            device_location = device_data["city"] + ", " + device_data["region"] + ", " + device_data["country"]

            # Create session
            device_data = {
                "device": device,
                "device_id": device_id,
                "device_location": device_location
            }
            url_path = "user/create-session"
            url = urljoin(backend_base_url, url_path)

            # print(device_data)
            response = requests.post(url, json=device_data)

            if not response.ok:
                print(f"[red]Unable to create session! Please try again later")
                return
            
            session_id = response.json()["data"][0]["session_id"]

            # Generater link & Open browser
            link = urljoin(frontend_base_url, f"verify-session/{session_id}")

            if browserless:
                print(f"Please open the following link in your browser to login: [underline]{link}[/underline]")
            else:
                # Open browser
                print(f"Opening the browser : [underline]{link}[/underline]")
                typer.launch(link)
            
            with Progress(
                SpinnerColumn(),
                TextColumn("[progress.description]{task.description}"),
                transient=True,
            ) as progress:
                progress.add_task(description="Waiting for response...", total=None)
                # Hit the endpoint to check if the session is verified
                start_time = time()
                while True:
                    url_path = "user/session-token"
                    url = urljoin(backend_base_url, url_path)
                    data = {
                        "session_id": session_id,
                        "device_id": device_id,
                    }
                    # print(data)
                    response = requests.post(url, json=data)
                    # print("Hit API response ", response.text)
                    if response.ok:
                        token = response.text
                        token = json.loads(token)["data"][0]

                        access_token = token["accessToken"]
                        email = token["email"]

                        save_auth(access_token=access_token, email=email)

                        break
                    else:
                        if response.status_code == 404:
                            print(f"[red]Session not found! Please try again later")
                            return
                        elif response.status_code == 403:
                            print(f"[red]Session expired or invalid device! Please try again later")
                            return

                    if time() - start_time > 60 * 10:
                        print(f"[red]Session timed out!")
                        return
                    
                    sleep(1)
            print(f"[green]Successfully logged in as {email}!")
            
            # Select organization
            select()
    except Exception as e:
        print(f"[red]Unable to contact the server! Please try again later")
        print(e)


@app.command()
def status():
    exists = get_token()
    if exists:
        print(f"[green]You are logged in!")


@app.command()
def logout():
    exists = get_token()
    if exists:
        delete_token()

if __name__ == "__main__":
    app()
