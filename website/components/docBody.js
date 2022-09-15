import React from "react";

export default function DocBody({ data, mdContent}) {
    return (
        <React.Fragment>
            <div className="m-auto text-3xl">{data}</div>
            <div className="m-auto text-2xl" dangerouslySetInnerHTML={{ __html: mdContent }}></div>
        </React.Fragment>
        
    )
}
