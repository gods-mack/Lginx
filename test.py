import asyncio
import aiohttp
from aiohttp import ClientSession

async def fetch_html(url: str, session: ClientSession, **kwargs) -> str:
    resp = await session.request(method="GET", url=url, **kwargs)
    resp.raise_for_status()
    return await resp.text()

async def make_requests(url: str, **kwargs) -> None:
    async with ClientSession() as session:
        tasks = []
        for i in range(1,200):
            tasks.append(
                fetch_html(url=url, session=session, **kwargs)
            )
        results = await asyncio.gather(*tasks)

        # do something with results

if __name__ == "__main__":
    asyncio.run(make_requests(url='http://localhost/home/'))