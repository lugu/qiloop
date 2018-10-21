unsorted to do list:
- meta:
    - add list of imports into the packageDeclaration
    - for each file extract the PackageDeclaration
    - for each package:
        - create a typeset per package.
        - for each pacakge declaration within the package:
            - add the types to the typeset
    - for each package declaration:
        - for each import:
            - add the imported package typeset
- meta: allow for unresolved meta object
- meta: generate proxy using type set information
- meta: generate enum types
- meta: generate proxy with enum types
- meta: generate proxy with object reference
- proxy: implements Disconnect
- server: implement service directory to validate the approach
- NOTES.md: machine id discussion
- NOTES.md: IDL section: syntax, ref to libqi doc and sample
- NOTES.md: sequence diagram with multiple process
- NODES.md: document the capability event type
- client: an API to make cancellable call (not the default way?)
- server: handle cancel message types
- tests: use passive service to complement proxy test
- session: register callback for service disconnection
- doc: restore the bootstrap discussion
- example: fuzzy tester
- example: firewall messaging level
- NOTES.md: document properties (semantic, messages)
- client: add properties
- service: hide generic object methods
- service: handle properties
- service: signals registration
- session: refactor session namespace with minimal deps

- failed to parse the following signatures:
    2018/10/21 15:17:23 failed to generate IDL of _SemanticEngine: failed to generate proxy object _SemanticEngine: method callback failed for makeKnowledge: failed to parse parms of makeKnowledge: failed to parse signature: (s(({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}fs)<AgentGrd,concepts,confidence,userId>i)<TextProcessingContext,author,recever,language>)
    2018/10/21 15:17:26 failed to generate IDL of ALAnimatedSpeech: failed to generate proxy object ALAnimatedSpeech: method callback failed for _sayKnowledge: failed to parse parms of _sayKnowledge: failed to parse signature: ((iis({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}f({ii})<Duration,timeInfos>)<TimeGrd,concepts,confidence,reference>(im)<Expression,type,exp>)<Knowledge,fromSource,fromLanguage,fromText,author,time,exp>(({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}fs)<AgentGrd,concepts,confidence,userId>i)<TextProcessingContext,author,recever,language>m)
    2018/10/21 15:17:32 failed to generate IDL of ALDialog: failed to generate proxy object ALDialog: signal callback failed for _heardKnowledge: failed to parse signal of _heardKnowledge: failed to parse signature: ((iis({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}f({ii})<Duration,timeInfos>)<TimeGrd,concepts,confidence,reference>(im)<Expression,type,exp>)<Knowledge,fromSource,fromLanguage,fromText,author,time,exp>)
    2018/10/21 15:17:32 failed to generate IDL of ALMood: failed to generate proxy object ALMood: method callback failed for currentPersonState: failed to parse return of currentPersonState: failed to parse signature: ((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>((ff)<BodyLanguageEase,level,confidence>)<BodyLanguageState,ease>(ff)<Smile,value,confidence>((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>)<Expressions,calm,anger,joy,sorrow,laughter,excitement,surprise>)<PersonState,valence,attention,bodyLanguageState,smile,expressions>
    2018/10/21 15:17:36 failed to generate IDL of ALRobotMood: failed to generate proxy object ALRobotMood: method callback failed for getFullState: failed to parse return of getFullState: failed to parse signature: ((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>)<RobotFullState,pleasure,excitement>
