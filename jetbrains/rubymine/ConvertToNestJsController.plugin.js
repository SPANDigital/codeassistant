importPackage(java.lang)
importPackage(java.util)
importPackage(com.intellij.ide.scratch)
importPackage(com.intellij.openapi.actionSystem)
importPackage(com.intellij.psi)

const actionId = 'ConvertToNestJs'

let project = ide.getProject()
let actionManager = ide.applicationService(ActionManager.class)
let editorActions = actionManager.getAction("EditorPopupMenu")


ide.registerAction(actionId, "Convert to NestJs controller", (e) => {
        //const time = new SimpleDateFormat("HH:mm").format(new Date())
        fileEditorManager = FileEditorManager.getInstance(project)

        currentDoc = fileEditorManager.getSelectedTextEditor().getDocument()
        psiFile = PsiDocumentManager.getInstance(ide.getProject()).getPsiFile(currentDoc)
        vFile = psiFile.getOriginalFile().getVirtualFile()
        path = vFile.getPath()

        cmd = "/opt/homebrew/bin/codeassistant rails2nextjs convert --railstype controller --nestjstype controller --src " + path + " --dest /dev/stdout"
        fullCmd = ["/bin/sh","-c", cmd]

        application = ide.applicationService(Application.class)

        var R = Java.extend(Java.type("java.lang.Runnable"));
        var r = new R(function(){

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
editorActions.addAction(action)

ide.addShortcut(actionId, "ctrl 1")



