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
                logo: element.querySelector('img').src,
                tvUrl: '',
                error: false
            });
      });
      
      return result;
  });

  for(let i=0;i<data.length;i++){
    let element = data[i];
    // console.log(element.url);
    try
    {
        await page.goto(element.url);
        
        element.logo = await page.evaluate(() => document.querySelector("meta[property='og:image']").content);
        
        let bodyHTML = await page.evaluate(() => document.body.innerHTML);
        
        let re = /(http.*?\.m3u8)/g;
        let myArray = re.exec(bodyHTML);
        element.tvUrl = myArray[1];
    }
    catch(err){
        element.error = true;
    }
  };

  console.log('#EXTM3U');
  for(let i=0;i<data.length;i++){
     let element = data[i];
     if(element.error === true)
        continue;
      console.log('');
      console.log('#EXTINF:-1 tvg-logo="' + element.logo + '",' + element.title);
      console.log(element.tvUrl);
  };
  await browser.close();
})();