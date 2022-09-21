import IzeNavbar from "./izeNavbar";

export default function Layout({ children }) {
    return (
        <>
            <IzeNavbar />
            {children}
        </>

    )
}
