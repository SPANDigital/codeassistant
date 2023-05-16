import * as React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import Link from '@mui/material/Link';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import NotificationsIcon from '@mui/icons-material/Notifications';
import Deposits from './Deposits';
import Orders from './Orders';
import {useEffect, useState} from "react";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import ListItemButton from "@mui/material/ListItemButton";
import LibraryIcon from '@mui/icons-material/LibraryBooks';
import Icon from '@mui/material/Icon';
import { BrowserRouter as Router, Routes, Route, Link as RouterLink} from "react-router-dom";
import Landing from "./landing";
import Library from "./library";
import Receive from "./receive";

import "./app.css";

function Copyright(props: any) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
    {'Copyright © '}
    <Link color="inherit" href="https://spandigital.com/">
        SPAN Digital
    </Link>{' '}
    {new Date().getFullYear()}
    {'.'}
    </Typography>
);
}


const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
}

const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        '& .MuiDrawer-paper': {
            position: 'relative',
            whiteSpace: 'nowrap',
            width: drawerWidth,
            transition: theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
            boxSizing: 'border-box',
            ...(!open && {
                overflowX: 'hidden',
                transition: theme.transitions.create('width', {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.leavingScreen,
                }),
                width: theme.spacing(7),
                [theme.breakpoints.up('sm')]: {
                    width: theme.spacing(9),
                },
            }),
        },
    }),
);

const mdTheme = createTheme();

interface DashboardContentProps {
    data: object
    error: string
}

function DashboardContent({ data }: DashboardContentProps) {
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };

    return (
        <ThemeProvider theme={mdTheme}>
        <Router>
        <Box sx={{ display: 'flex' }}>
    <CssBaseline />
    <AppBar position="absolute" open={open}>
    <Toolbar
        sx={{
        pr: '24px', // keep right padding when drawer closed
    }}
>
    <IconButton
        edge="start"
    color="inherit"
    aria-label="open drawer"
    onClick={toggleDrawer}
    sx={{
        marginRight: '36px',
    ...(open && { display: 'none' }),
    }}
>
    <MenuIcon />
    </IconButton>
    <Typography
    component="h1"
    variant="h6"
    color="inherit"
    noWrap
    sx={{ flexGrow: 1 }}
>
    Dashboard
    </Typography>
    <IconButton color="inherit">
    <Badge badgeContent={4} color="secondary">
        <NotificationsIcon />
        </Badge>
        </IconButton>
        </Toolbar>
        </AppBar>
        <Drawer variant="permanent" open={open}>
    <Toolbar
        sx={{
        display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-end',
            px: [1],
    }}
>
    <IconButton onClick={toggleDrawer}>
    <ChevronLeftIcon />
    </IconButton>
    </Toolbar>
    <Divider />
    <List component="nav">
        { data && Object.values(data).map((value, index) => {
            return (
                <ListItemButton component={RouterLink} to={ "/web/" + value.Name}>
                    <ListItemIcon>
                        <Icon>{value.Icon.toLowerCase()}</Icon>
                    </ListItemIcon>
                    <ListItemText primary={ value.DisplayName } />
                </ListItemButton>
            );
        })}
        <Divider sx={{ my: 1 }} />
    </List>
    </Drawer>
    <Box
    component="main"
    sx={{
        backgroundColor: (theme) =>
            theme.palette.mode === 'light'
                ? theme.palette.grey[100]
                : theme.palette.grey[900],
            flexGrow: 1,
            height: '100vh',
            overflow: 'auto',
    }}
>
    <Toolbar />
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Routes>
           <Route path="/web" element={<Landing />} />
           <Route path="/web/:libraryName" element={<Library data={data} />} />
           <Route path="/web/:libraryName/:commandName/:receiveChannel" element={<Receive data={data}  />} />
        </Routes>
    <Copyright sx={{ pt: 4 }} />
    </Container>
    </Box>
    </Box>
    </Router>
    </ThemeProvider>
);
}

export default function App() {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch(`/api/graph`)
            .then((response) => {
                if (!response.ok) {
                    throw new Error(
                        `This is an HTTP error: The status is ${response.status}`
                    );
                }
                return response.json();
            })
            .then((actualData) => {
                setData(actualData);
                setError(null);
            })
            .catch((err) => {
                setError(err.message);
                setData(null);
            })
            .finally(() => {
                setLoading(false);
            });
    }, []);

    return <DashboardContent data={data} error={error} />;
}