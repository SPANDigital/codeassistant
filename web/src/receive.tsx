import {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Paper from "@mui/material/Paper";
import ReactMarkdown from "react-markdown";
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {dark} from 'react-syntax-highlighter/dist/esm/styles/prism'
import CopyToClipboard from 'react-copy-to-clipboard';

interface ReceiveContentProps {
    data: object
}

export default function Receive({ data }: ReceiveContentProps) {
    const [message, setMessage] = useState("")

    let {libraryName, commandName, receiveChannel} = useParams()

    let library = data[libraryName]
    let command = library.Commands[commandName]

    useEffect(() => {
        const eventSource = new EventSource("/api/receive/" + receiveChannel)
        eventSource.onmessage = (event) => {
            let eventData = JSON.parse(event.data)
            if (eventData.Type == "Part") {
                setMessage((message) => message + eventData.Delta)
            }
        }
        return () => {
            eventSource.close();
        };
    }, [])

    return (
        <React.Fragment>
            <Grid item xs={12}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <Typography component="h2" variant="h5" color="primary" gutterBottom>
                        { library.DisplayName } &gt; { command.DisplayName }
                    </Typography>
                    { message != "" &&
                        <ReactMarkdown children={message} components={{
                            code({node, inline, className, children, ...props}) {
                                const match = /language-(\w+)/.exec(className || '')
                                return !inline && match ? (
                                    <React.Fragment>
                                        <CopyToClipboard text={String(children).replace(/\n$/, '')}>
                                            <button>Copy</button>
                                        </CopyToClipboard>
                                        <SyntaxHighlighter
                                            {...props}
                                            children={String(children).replace(/\n$/, '')}
                                            style={dark}
                                            language={match[1]}
                                            PreTag="div"
                                        />
                                    </React.Fragment>
                                ) : (
                                    <code {...props} className={className}>
                                        {children}
                                    </code>
                                )
                            }
                        }}
                        />
                    }
                </Paper>
            </Grid>
        </React.Fragment>
    );
}