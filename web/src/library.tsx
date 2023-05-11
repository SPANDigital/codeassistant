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
import {Button, FormControl, Input, InputLabel, ListItem, TextField} from "@mui/material";

interface LibraryContentProps {
    data: object
}

export default function Library({ data }: LibraryContentProps) {
    let { libraryName } = useParams();
    let library = data[libraryName]
    return (
        <React.Fragment>
            <Grid item xs={12}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <Typography component="h2" variant="h6" color="primary" gutterBottom>
                        { library.DisplayName }
                    </Typography>

                        { data && Object.values(library.Commands).map((value, index) => {
                            return (
                                <React.Fragment>
                                    <Typography component="h3" variant="h7" color="secondary" gutterBottom>
                                        { value.Name }
                                    </Typography>
                                        { value.Params && Object.keys(value.Params).map((key, index) => {
                                           return (
                                               <FormControl>
                                                <InputLabel htmlFor={key}>{ key }</InputLabel>
                                                <Input id={ key } />
                                               </FormControl>
                                           )
                                        })}
                                    <Button>Fetch</Button>
                                    <Divider sx={{ my: 1 }} />
                                </React.Fragment>
                            );
                        })}


                </Paper>
            </Grid>
        </React.Fragment>
    );
}