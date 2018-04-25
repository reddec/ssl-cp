from jinja2 import Environment, FileSystemLoader
import os


def render(*asset_path, **kwargs):
    env = Environment(loader=FileSystemLoader(os.path.join(os.path.dirname(__file__), '..', '..', "assets")))
    templ = env.get_template(os.path.join(*asset_path))
    return templ.render(**kwargs).encode()


from . import stunnel, nginx
