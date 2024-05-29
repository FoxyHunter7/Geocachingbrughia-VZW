<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class EventLanguage extends Model
{
    use HasFactory;

    protected $fillable = [
        "event_id",
        "lang_code",
        "description"
    ];

    protected $hidden = [
        'id',
        'event_id',
        'created_at',
        'updated_at',
    ];
}
