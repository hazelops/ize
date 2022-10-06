import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

import 'tailwindcss/tailwind.css'
import '../main.css'

import Layout from '../components//layouts/layout'

function MyApp({ Component, pageProps }) {
  library.add(fas)
  return (
    <Layout>
      <Component {...pageProps} />
    </Layout>
  ) 
}

export default MyApp
