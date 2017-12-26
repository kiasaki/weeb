import { bootstrap } from './framework';

import container from './container';

bootstrap(container);

const app = container.get("app");
app.start();


