"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const images_1 = require("./images");
const names_1 = require("./names");
const app = (0, express_1.default)();
const port = 3000;
// Fake data creation function
const data = {
    create: () => ({
        repositories: ['busybox', 'centos', 'hello-world'],
        tags: {
            busybox: ['latest', '1.33.1'],
            centos: ['latest', '7', '8'],
            'hello-world': ['latest'],
        },
    }),
};
const mockRepos = {};
const randomizeImageName = (imageName) => {
    const randomString = Math.random().toString(36).substring(7);
    const randomName = names_1.names[Math.floor(Math.random() * names_1.names.length)];
    const parts = imageName.split('/');
    parts.push(names_1.names[Math.floor(Math.random() * names_1.names.length)]);
    return [randomName, parts[0], parts[1] + '-' + randomString].join('/');
};
const createRandomTags = (max = 10) => {
    const res = [];
    const random = Math.floor(Math.random() * max);
    for (let i = 0; i < random; i++) {
        res.push(createRandomSemver());
    }
    return res;
};
function delay(t = 300) {
    return new Promise((resolve) => setTimeout(resolve, t));
}
const createRandomImage = (max = 1000) => {
    const res = [];
    const random = Math.floor(Math.random() * max);
    for (let i = 0; i < random; i++) {
        const randomImage = images_1.images[Math.floor(Math.random() * images_1.images.length)];
        if (Math.random() > 0.7) {
            res.push(randomImage);
        }
        else {
            res.push(randomizeImageName(randomImage));
        }
    }
    return res;
};
const createRandomSemver = () => {
    const major = Math.floor(Math.random() * 10);
    const minor = Math.floor(Math.random() * 10);
    const patch = Math.floor(Math.random() * 10);
    return `${major}.${minor}.${patch}`;
};
// Endpoint to list Docker repositories
app.get('/api/docker/:repoKey/v2/_catalog', (req, res) => {
    console.log(req.url);
    const { repoKey } = req.params;
    const result = createRandomImage(500);
    delay().then(() => {
        res.json({ repositories: result });
    });
});
// Endpoint to list Docker tags
app.get('/api/docker/:repoKey/v2/:imageName*/tags/list', (req, res) => {
    console.log(req.url);
    const { repoKey, imageName } = req.params;
    const result = data.create();
    delay().then(() => {
        res.json({ name: imageName, tags: createRandomTags() });
    });
    // res.json({ name: imageName, tags: createRandomTags() });
});
app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});
