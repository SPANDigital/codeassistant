import Typography from "@mui/material/Typography";
import * as React from "react";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid";
import {Link as RouterLink, useParams} from "react-router-dom";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import Icon from "@mui/material/Icon";
import ListItemText from "@mui/material/ListItemText";
import Divider from "@mui/material/Divider";
import List from "@mui/material/List";
import {Button, FormControl, FormHelperText, Input, InputLabel, ListItem, TextField} from "@mui/material";
import ReactMarkdown from "react-markdown";
import {useState} from "react";
import PromptResponse from "./promptresponse";

interface LibraryContentProps {
    data: object
}

export default function Library({ data }: LibraryContentProps) {
    const [activeCommand, setActiveCommand] = useState("")
    const [message, setMessage] = useState("")


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

                        { data && Object.values(library.Commands).filter(command => activeCommand == "" || command.Name == activeCommand).map((command, index) => {
                            let handleSubmit = (event) => {
                                event.preventDefault()
                                setActiveCommand(command.Name)
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
                                        console.log("Event source location", location)
                                        let eventSource = new EventSource(location)
                                        let message = ""
                                        eventSource.onmessage = (event) => {
                                            let eventData = JSON.parse(event.data)
                                            if (eventData.Type == "Part") {
                                                message = message + eventData.Delta
                                                setMessage(message)
                                            }
                                        }
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