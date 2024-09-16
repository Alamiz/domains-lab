import { useState } from "react";
import { RxHamburgerMenu, RxCross2 } from "react-icons/rx";

const NavBar = () => {
  const [currentTab, setCurrentTab] = useState('home');
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="p-6 bg-background/85 backdrop-blur sticky top-0 inset-x-0 z-40 shadow-md">
      <div className="container">
        <nav className="flex items-center justify-between gap-6">
          <a className="w-52 cursor-pointer" href="#">
            <img src="./images/Domain-Lab-Logo.svg" alt="" />
          </a>
          <ul className="flex flex-1 items-center justify-center gap-4 md:flex hidden">
            <li onClick={() => setCurrentTab('home')}><a href="#" className={currentTab === "home" ? "text-primary underline underline-offset-8" : "text-secondary"}>Home</a></li>
            <li onClick={() => setCurrentTab('upload')}><a href="#upload" className={currentTab === "upload" ? "text-primary underline underline-offset-8" : "text-secondary"}>Upload File</a></li>
            <li onClick={() => setCurrentTab('search')}><a href="#search" className={currentTab === "search" ? "text-primary underline underline-offset-8" : "text-secondary"}>Search</a></li>
            <li onClick={() => setCurrentTab('download')}><a href="#download" className={currentTab === "download" ? "text-primary underline underline-offset-8" : "text-secondary"}>Download</a></li>
          </ul>
          <button className="md:flex hidden bg-primary rounded-md px-3 py-2 text-white">Get started</button>
          <button className="md:hidden" onClick={() => setIsOpen(!isOpen)}>
            {isOpen ? <RxCross2 />: <RxHamburgerMenu />}
          </button>

        </nav>
          {isOpen && (
            <ul className="flex flex-1 flex-col mt-4 items-center justify-center gap-4">
              <li onClick={() => setIsOpen(false)}><a href="#" className={"text-secondary"}>Home</a></li>
              <li onClick={() => setIsOpen(false)}><a href="#upload" className="text-secondary">Upload File</a></li>
              <li onClick={() => setIsOpen(false)}><a href="#search" className="text-secondary">Search</a></li>
              <li onClick={() => setIsOpen(false)}><a href="#download" className="text-secondary">Download</a></li>
            </ul>
          )}
      </div>
    </header>
  )
}

export default NavBar