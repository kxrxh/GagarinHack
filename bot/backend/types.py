from dataclasses import dataclass, field


@dataclass
class HumanInfo:
    name: str
    sex: str
    birth_date: str
    death_date: str

    # Optional
    question_answers: dict[str, str] = field(default_factory=dict)

    # Generated
    epitaph: str = field(default="")
    biography: str = field(default="")


@dataclass
class GenerationRequest:
    name: str
    sex: str
    birth_date: str
    death_date: str
    question_answers: dict[str, str] = field(default_factory=dict)
