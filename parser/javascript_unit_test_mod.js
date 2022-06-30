"use strict";

const {promises: fs} = require('fs');
const app = express();

const router = express.Router();
router.get('/', async (req, res) => {
    const f = await fs.readFile("/tmp/somefile.txt");
    res.send(f.toString("utf-8"));
});
app.use(router);

const f2 = await fs.readFile("/tmp/somefile_2.txt");
console.log(f2);

app.listen(3000, async () => {
    console.log('App listening locally at :3000');
});