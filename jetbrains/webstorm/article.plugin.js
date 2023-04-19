importPackage(java.lang)
importPackage(java.util)
importPackage(com.intellij.ide.scratch)
importPackage(com.intellij.openapi.command)
importPackage(com.intellij.openapi.actionSystem)
importPackage(com.intellij.psi)

const actionId = 'RetrieveArticleFromChatGPT'

let project = ide.getProject()
let actionManager = ide.applicationService(ActionManager.class)
let editorActions = actionManager.getAction("EditorPopupMenu")
let application = ide.applicationService(Application.class)

ide.registerAction(actionId, "Lookup article using ChatGPT", (e) => {

    let fileEditorManager = FileEditorManager.getInstance(project)
    let editor = fileEditorManager.getSelectedTextEditor()
    let currentDoc = editor.getDocument()
    let primaryCaret = editor.getCaretModel().getPrimaryCaret();
    let start = primaryCaret.getSelectionStart();
    let end = primaryCaret.getSelectionEnd();

    let term = currentDoc.getText().subSequence(start, end)
    primaryCaret.removeSelection()

    ide.notify("Article retriever", "Writing article about " + term)

    let cmd = "/opt/homebrew/bin/codeassistant article \"" + term + "\""
    let fullCmd = ["/bin/sh","-c", cmd]

        var R = Java.extend(Java.type("java.lang.Runnable"));
        var r = new R(function(){


            process =new ProcessBuilder(fullCmd).start()

            scanner = new Scanner(process.getInputStream()).useDelimiter("\\A");
            result = scanner.hasNext() ? scanner.next() : "";

            ide.notify("Article retriever", "Gpt result: " + result)

            scratchFile = ScratchRootType.getInstance().createScratchFile(project, term +".md", null, result, ScratchFileService.Option.create_new_always)

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



