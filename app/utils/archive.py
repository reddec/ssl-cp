import tarfile
import io
from typing import Dict


def add_to_archive(tar: tarfile.TarFile, buf: bytes, name: str):
    info = tarfile.TarInfo(name)
    info.size = len(buf)
    info.type = tarfile.REGTYPE
    tar.addfile(info, io.BytesIO(buf))


def dict_to_archive(files: Dict[str, bytes]) -> bytes:
    mem_arch = io.BytesIO()
    with tarfile.open(fileobj=mem_arch, mode='w:gz') as tar:
        for filename, content in files.items():
            add_to_archive(tar, content, filename)
    mem_arch.seek(0)
    return mem_arch.read()
