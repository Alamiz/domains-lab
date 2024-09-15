import { Header, FileUpload, FileDownload, HowItWorks, Footer } from './containers'
import { NavBar, CTA } from './components'
import { useState } from 'react'

const App = () => {

  return (
    <div className="app">
      <NavBar />
      <Header />
      <FileUpload />
      <FileDownload />
      <HowItWorks />
      <Footer />
    </div>
  )
}

export default App