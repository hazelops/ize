import 'tailwindcss/tailwind.css'
import '../main.css'

import Layout from '../components/layout'

function MyApp({ Component, pageProps }) {
  return (
    <Layout>
      <Component {...pageProps} />
    </Layout>
  ) 
}

export default MyApp
