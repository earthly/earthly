# Next.js on Netlify Example

This page will walk you through an example of how to build a Next.js hello world website and deploy it to Netlify using Earthly.  
  
For this example, we are using the [basic-css](https://github.com/vercel/next.js/tree/canary/examples/basic-css) Next.js example that is included in [Vercel's Next.js repository](https://github.com/vercel/next.js).  
  
This example has two files that contain all of the code and CSS for this site.  
  
`pages\index.tsx`
```
import styles from '../styles.module.css'

const Home = () => {
  return (
    <div className={styles.hello}>
      <p>Hello World</p>
    </div>
  )
}

export default Home
```
  
`styles.module.css`
```
.hello {
  font: 15px Helvetica, Arial, sans-serif;
  background: #eee;
  padding: 100px;
  text-align: center;
  transition: 100ms ease-in background;
}
.hello:hover {
  background: #ccc;
}
```
The standard way to build and deploy is to build locally with `npm run build` to ensure that the site builds and then push the changes into your `main` git branch and let Netlify detect and auto-deploy the changes. There are a few limitations to this standard approach though:
1.  It is not possible to fully verify that your build and deploy will both complete successfully in Netlify. Netlify uses their own libraries to build and deploy that function slightly differently than `npm run build`.
2.  You have to grant Netlify access to your repository. Generally, this is not a big issue, but it can be more difficult if your repository is under an organization you don't own or if you want to use a private repository.
  
Instead, we will be using an Earthfile to build and deploy the site to Netlify. We will create three targets: `deps`, `build`, and `deploy`.
1.  `deps` is the base target that `build` and `deploy` use. It copies all of the common files and sets environment variables that both the `build` and `deploy` targets require.
2.  `build` is a target that copies the local files and directories required for building, installs prerequisites, and builds your website using Netlify's CLI.
3.  `deploy` is a target that runs the `build` target (if it isn't cached), copies the output from `build` that is required for deploying to Netlify, and deploys your site to Netlify using their CLI.
  
**Note:** We specify the platform in the base image to be linux/amd64. This is because the Netlify CLI uses [Deno](https://deno.land/), and Deno doesn't support linux/arm64 at the moment. If you're using an arm64 computer, the `build` target will fail unless the platform is specified as amd64.  
  
`Earthfile`
```
VERSION 0.6
FROM --platform=linux/amd64 node:latest
WORKDIR /app

deps:
    ARG auth_token
    ARG site_id
    # Set environment variables required for building and deploying to Netlify
    ENV NETLIFY_AUTH_TOKEN=$auth_token
    ENV NETLIFY_SITE_ID=$site_id
    # Install the Netlify CLI (global)
    RUN npm install netlify-cli -g
    # Copy package.json for installing required packages
    COPY package.json ./
    # Install the Netlify CLI (local) and Netlify Next.js plugin and required packages
    RUN npm install netlify-cli --save-dev
    RUN npm install @netlify/plugin-nextjs --save-dev
    RUN npm install
    # Copy netlify.toml - required for building and deploying to Netlify
    COPY netlify.toml ./

build:
    FROM +deps
    # Copy files and directories required for building
    COPY next-env.d.ts styles.module.css tsconfig.json ./
    COPY --dir pages ./
    # Build site
    RUN netlify build --context production
    SAVE ARTIFACT ./node_modules node_modules/ AS LOCAL ./
    SAVE ARTIFACT ./.next .next/ AS LOCAL ./
    SAVE ARTIFACT ./.netlify .netlify/ AS LOCAL ./

deploy:
    FROM +deps
    # Copy artifacts required for deploying to Netlify
    COPY +build/node_modules/ +build/.next/ +build/.netlify/ ./
    # Deploy site
    RUN netlify deploy --prod
```
We also have to create a `netlify.toml` file.
```
[[plugins]]
package = "@netlify/plugin-nextjs"

[build]
command = "npm run build"
publish = ".next"

[[redirects]]
from = "/_next/static/*"
to = "/static/:splat"
status = 200
force = true
```
We must supply two arguments for `build` and `deploy` to function properly. Read [our docs](https://docs.earthly.dev/docs/guides/build-args) for more information on the different ways to supply arguments for an Earthly build.
1.  `auth_token`:  Sets the `NETLIFY_AUTH_TOKEN` environment variable. Read [Netlify's docs](https://docs.netlify.com/cli/get-started/#authentication) for more information on how to get an auth token.
2.  `site_id`:  Sets the `NETLIFY_SITE_ID` environment variable. Read [Netlify's docs](https://docs.netlify.com/cli/get-started/#link-with-an-environment-variable) for more information on how to get your site id.
  
To verify that your site builds correctly but not deploy, run `earthly +build` (assumes using a .env file to supply arguments) or `earthly +build --auth_token=[your_auth_token] --site_id=[your_site_id]`.
  
To build and deploy your site, run `earthly +deploy` (assumes using a .env file to supply arguments) or `earthly +deploy --auth_token=[your_auth_token] --site_id=[your_site_id]`.
