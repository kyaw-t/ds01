let count = 100;
const fs = require('fs');
async function getImages(start) {
  fetch(
    `https://hub.docker.com/api/search/v3/catalog/search?query=&from=${start}&open_source=false&official=false&size=200`
  )
    .then((response) => response.json())
    .then(async (data) => {
      for (d in data.results) {
        const image = data.results[d];
        console.log(`${JSON.stringify(image.name)},`);
        count += 1;
        // images.push(image);
      }
      // console.log('_', images.length);
      if (count < 7000) {
        await getImages(count);
      }
    });
}

getImages(count);
