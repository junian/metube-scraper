const puppeteer = require('puppeteer');
const metubeLiveUrl = 'https://www.metube.id/live';

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto(metubeLiveUrl);
  
  let data = await page.evaluate(() => 
  {
      let thumbnails = document.querySelectorAll('.livetv-thumbnail');
      let result = [];
      thumbnails.forEach(element => {
        result.push(
            {
                title: element.querySelector("a").title,
                url: element.querySelector("a").href,
                logo: element.querySelector('img').src
            });
      });
      
      return result;
  });
  console.log(data);
  await browser.close();
})();