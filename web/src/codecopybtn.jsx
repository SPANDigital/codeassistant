import * as React from "react";

export default function CodeCopyBtn({ children }) {
    const [copyOk, setCopyOk] = React.useState(false);
    const iconColor = copyOk ? '#0af20a' : '#ddd';
    const icon = copyOk ? 'fa-check-square' : 'fa-copy';
    const handleClick = (e) => {
        navigator.clipboard.writeText(children[0].props.children[0]);
        console.log(children)
        setCopyOk(true);
        setTimeout(() => {
            setCopyOk(false);
        }, 500);
    }

    return  (
        <div className="code-copy-btn">
            <i className={`fas ${icon}`} onClick={handleClick}></i>
        </div>
    );
}
