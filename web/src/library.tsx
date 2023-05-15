import Typography from "@mui/material/Typography";
import * as React from "react";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid";
import {useParams} from "react-router-dom";
import Divider from "@mui/material/Divider";
import {Button,TextField} from "@mui/material";
import ReactMarkdown from "react-markdown";
import {useState} from "react";
import { useNavigate } from "react-router-dom";

interface LibraryContentProps {
    data: object
}

export default function Library({ data }: LibraryContentProps) {
    const [message, setMessage] = useState("")

    const navigate = useNavigate();

    let { libraryName } = useParams();
    let library = data[libraryName]

    return (
        <React.Fragment>
            <Grid item xs={12}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <Typography component="h2" variant="h5" color="primary" gutterBottom>
                        { library.DisplayName }
                    </Typography>

                    { (library.Index && library.Index.trim() != "") &&
                        <ReactMarkdown>{library.Index}</ReactMarkdown>
                    }

                        { data && Object.values(library.Commands).filter((command) => !command.Abstract).map((command, index) => {
                            let handleSubmit = (event) => {
                                event.preventDefault()
                                let target   = event.target;
                                const formData  = new FormData();
                                for (let i = 0; i < target.length; i++) {
                                    let name = target.elements[i].getAttribute("name")
                                    if (name != null) {
                                        formData.append(name, target.elements[i].value)
                                    }
                                }
                                fetch(`/api/prompt/` + library.Name + '/' + command.Name, {
                                    method: 'POST',
                                    body: formData
                                })
                                    .then((response) => {
                                        if (!response.ok) {
                                            throw new Error(
                                                `This is an HTTP error: The status is ${response.status}`
                                            );
                                        }
                                        return response.headers.get("Location")
                                    })
                                    .then((location) => {
                                        let channel = location.split('/').slice(-1)[0]
                                        console.log('Found channel', channel)
                                        navigate('/web/' + library.Name + '/' + command.Name + '/' + channel)
                                    })
                                    .catch((err) => {
                                        console.log(err)
                                    })
                            }
                            return (
                                <React.Fragment>
                                    <Typography component="h3" variant="h6" color="secondary" gutterBottom>
                                        { command.DisplayName }
                                    </Typography>
                                        <form onSubmit={ handleSubmit }>
                                        { command.Params && Object.keys(command.Params).map((key, index) => {
                                            console.log('value', command)
                                            let uiHints = command.UiHints || {}
                                            let uiHint = {...{
                                                Label: key,
                                                HelperText: key,
                                                Props: {}
                                            }, ...(uiHints[key] || {})}

                                           return (
                                               <TextField id={ key } name={ key } label={ uiHint.label }  helperText={ uiHint.HelperText } {...uiHint.Props} fullWidth />
                                           )
                                        })}
                                            <Button type="submit">Fetch</Button>
                                        </form>

                                    <Divider sx={{ my: 1 }} />
                                    { message != "" &&
                                        <ReactMarkdown>{ message }</ReactMarkdown>
                                    }
                                </React.Fragment>
                            );
                        })}


                </Paper>
            </Grid>
        </React.Fragment>
    );
}