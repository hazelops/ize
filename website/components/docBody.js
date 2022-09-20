import React from "react";

export default function DocBody({ mdContent }) {
    return (
        <React.Fragment>
            <div className="m-auto text-2xl" dangerouslySetInnerHTML={{ __html: mdContent }}></div>
        </React.Fragment>
        
    )
}
