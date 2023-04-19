importPackage(java.lang)
importPackage(java.util)
importPackage(com.intellij.ide.scratch)
importPackage(com.intellij.openapi.actionSystem)
importPackage(com.intellij.psi)

const actionId = 'ConvertToNestJsJestTest'

let project = ide.getProject()
let actionManager = ide.applicationService(ActionManager.class)
let editorActions = actionManager.getAction("EditorPopupMenu")

ide.registerAction(actionId, "Convert to NestJs Jest Test", (e) => {
        //const time = new SimpleDateFormat("HH:mm").format(new Date())
        project = ide.getProject()
    let fileEditorManager = FileEditorManager.getInstance(project)

    let currentDoc = fileEditorManager.getSelectedTextEditor().getDocument()
        psiFile = PsiDocumentManager.getInstance(ide.getProject()).getPsiFile(currentDoc)
        vFile = psiFile.getOriginalFile().getVirtualFile()
        path = vFile.getPath()

        cmd = "/opt/homebrew/bin/codeassistant rails2nextjs convert --railstype spec --nestjstype \"test using jest\" --src " + path + " --dest /dev/stdout"
        fullCmd = ["/bin/sh","-c", cmd]

        application = ide.applicationService(Application.class)

        const R = Java.extend(Java.type("java.lang.Runnable"));
        let r = new R(function(){

            process =new ProcessBuilder(fullCmd).start()
            ide.notify("Code generator", "Generated generation started")

            scanner = new Scanner(process.getInputStream()).useDelimiter("\\A");
            result = scanner.hasNext() ? scanner.next() : "";

            scratchFile = ScratchRootType.getInstance().createScratchFile(project, "controller.ts", null, result, ScratchFileService.Option.create_new_always)

            ide.notify("Code generated", "Generated generation complete")

            var openEditor = new R(function() {
                fileEditorManager.openTextEditor(new OpenFileDescriptor(project, scratchFile), true)
            })
            application.invokeLater(openEditor)


        });

        application.executeOnPooledThread(r)

    }
)

action = actionManager.getAction(actionId)
if (!editorActions.containsAction(action)) editorActions.addAction(action)

ide.addShortcut(actionId, "ctrl 2")

