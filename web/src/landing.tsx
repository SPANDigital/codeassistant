import Typography from "@mui/material/Typography";
import * as React from "react";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid";

export default function Landing() {
    return (
        <React.Fragment>
            <Grid item xs={12}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <Typography component="h2" variant="h6" color="primary" gutterBottom>
                        Landing page
                    </Typography>
                </Paper>
            </Grid>
        </React.Fragment>
    );
}